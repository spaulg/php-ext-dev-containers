ifneq ($(DEB_HOST_ARCH),$(filter $(DEB_HOST_ARCH),hurd-i386 m68k hppa ppc64))
  ext_PACKAGES += interbase
  interbase_DESCRIPTION := Interbase
  interbase_EXTENSIONS  := pdo_firebird
  pdo_firebird_config   := --with-pdo-firebird=shared,/usr
  export interbase_EXTENSIONS
  export interbase_DESCRIPTION
endif
