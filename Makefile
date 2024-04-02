build:
	hack/build.sh $(OS) $(ARCH)

package:
	hack/package.sh $(OS) $(ARCH)

LICENSE_CODEGEN=addlicense
LICENSE_DIRS=cmd hack pkg
export OPENAPI_PATH=

install-addlicense:
	go install github.com/google/addlicense@v1.1.1

generate-license:
	@for DIR in $(LICENSE_DIRS) ; do \
		echo "patched license in $$DIR"; \
		$(LICENSE_CODEGEN) -c "The mobaxterm-keygen Authors." -y "2024" -l "Apache-2.0" $$DIR; \
	done
