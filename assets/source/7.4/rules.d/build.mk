override_dh_auto_configure-indep:

override_dh_auto_configure-arch: prepared
	for target in $(TARGETS); do \
	  dh_auto_configure --builddirectory $${target}-build $(PARALLEL) -- $$(eval echo \$${$${target}_config}); \
	done

override_dh_auto_build-indep:

override_dh_auto_build-arch:
	for target in $(TARGETS); do \
	  dh_auto_build --builddirectory $${target}-build $(PARALLEL); \
	done

override_dh_auto_test-indep:

override_dh_auto_test-arch:
