package cmd

import (
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/artifact"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"
)

type task func(*sbom.Artifacts, *source.Source) ([]artifact.Relationship, error)

func catalogingTasks() ([]task, error) {
	var tasks []task

	generators := []func() (task, error){
		generateCatalogPackagesTask,
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
