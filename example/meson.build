example_moc_files = qt5.preprocess(
  moc_headers: 'example.h',
)

exe = executable('fimbeb-example',
                 'example.cpp',
                 example_moc_files,
                 link_with: fimbeb_lib,
                 dependencies: qt5_deps,
                 include_directories: fimbeb_headers)