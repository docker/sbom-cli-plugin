package cmd

import (
	"crypto"
	"fmt"

	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/artifact"
	"github.com/anchore/syft/syft/file"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"
)

type task func(*sbom.Artifacts, *source.Source) ([]artifact.Relationship, error)

func catalogingTasks() ([]task, error) {
	var tasks []task

	generators := []func() (task, error){
		generateCatalogPackagesTask,
		generateCatalogFileMetadataTask,
		generateCatalogFileDigestsTask,
	}

	for _, generator := range generators {
		task, err := generator()
		if err != nil {
			return nil, err
		}

		if task != nil {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func generateCatalogPackagesTask() (task, error) {
	if !appConfig.Package.Cataloger.Enabled {
		return nil, nil
	}

	task := func(results *sbom.Artifacts, src *source.Source) ([]artifact.Relationship, error) {
		packageCatalog, relationships, theDistro, err := syft.CatalogPackages(src, appConfig.Package.ToConfig())
		if err != nil {
			return nil, err
		}

		results.PackageCatalog = packageCatalog
		results.LinuxDistribution = theDistro

		return relationships, nil
	}

	return task, nil
}

func generateCatalogFileMetadataTask() (task, error) {
	if !appConfig.FileMetadata.Cataloger.Enabled {
		return nil, nil
	}

	metadataCataloger := file.NewMetadataCataloger()

	task := func(results *sbom.Artifacts, src *source.Source) ([]artifact.Relationship, error) {
		resolver, err := src.FileResolver(appConfig.FileMetadata.Cataloger.ScopeOpt)
		if err != nil {
			return nil, err
		}

		result, err := metadataCataloger.Catalog(resolver)
		if err != nil {
			return nil, err
		}
		results.FileMetadata = result
		return nil, nil
	}

	return task, nil
}

func generateCatalogFileDigestsTask() (task, error) {
	if !appConfig.FileMetadata.Cataloger.Enabled {
		return nil, nil
	}

	supportedHashAlgorithms := make(map[string]crypto.Hash)
	for _, h := range []crypto.Hash{
		crypto.MD5,
		crypto.SHA1,
		crypto.SHA256,
	} {
		supportedHashAlgorithms[file.DigestAlgorithmName(h)] = h
	}

	var hashes []crypto.Hash
	for _, hashStr := range appConfig.FileMetadata.Digests {
		name := file.CleanDigestAlgorithmName(hashStr)
		hashObj, ok := supportedHashAlgorithms[name]
		if !ok {
			return nil, fmt.Errorf("unsupported hash algorithm: %s", hashStr)
		}
		hashes = append(hashes, hashObj)
	}

	digestsCataloger, err := file.NewDigestsCataloger(hashes)
	if err != nil {
		return nil, err
	}

	task := func(results *sbom.Artifacts, src *source.Source) ([]artifact.Relationship, error) {
		resolver, err := src.FileResolver(appConfig.FileMetadata.Cataloger.ScopeOpt)
		if err != nil {
			return nil, err
		}

		result, err := digestsCataloger.Catalog(resolver)
		if err != nil {
			return nil, err
		}
		results.FileDigests = result
		return nil, nil
	}

	return task, nil
}

func runTask(t task, a *sbom.Artifacts, src *source.Source, c chan<- artifact.Relationship, errs chan<- error) {
	defer close(c)

	relationships, err := t(a, src)
	if err != nil {
		errs <- err
		return
	}

	for _, relationship := range relationships {
		c <- relationship
	}
}
