


asdf:
	asdf plugin add golang; if [ $$? -eq 2 ] ; then true ; fi
	asdf plugin add nodejs; if [ $$? -eq 2 ] ; then true ; fi
	asdf plugin add protoc; if [ $$? -eq 2 ] ; then true ; fi
	asdf plugin add goreleaser; if [ $$? -eq 2 ] ; then true ; fi
	asdf plugin add golangci-lint; if [ $$? -eq 2 ] ; then true ; fi
	asdf install
