package subroutines

//go:generate go run -mod=mod github.com/vektra/mockery/v2@v2.46.0 --all --case=underscore --with-expecter
//go:generate go run -mod=mod github.com/vektra/mockery/v2@v2.46.0 --srcpkg=sigs.k8s.io/controller-runtime/pkg/client --name=Client --case=underscore --with-expecter
