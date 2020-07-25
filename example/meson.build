example_moc_files = qt5.preprocess(
  moc_headers: 'example.gen.h',
)

exe = executable('fimbeb-example',
                 'example.cpp',
                 example_moc_files,
                 link_with: pokipoki_lib,
                 dependencies: qt5_deps,
                 include_directories: pokipoki_headers)