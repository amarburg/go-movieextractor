
SUBDIRS = movieextractor multimov framesource frameset virtualmov

test:
	for subdir in ${SUBDIRS} ; do \
		cd $$subdir && go test ; cd .. ; \
  done


.PHONY: test
