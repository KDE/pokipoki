src = [
  'Database.cpp',
]

headers = [
  'headers/Database.h',
]

moc_files = qt5.preprocess(
  moc_headers: headers,
)

install_headers(headers)

pokipoki_lib = shared_library('pokipoki',
                             src,
                             headers,
                             moc_files,
                             dependencies: [
                                 qt5_deps,
                             ],
                             include_directories: pokipoki_headers,
                             install: true)

pkgconfig = import('pkgconfig')
pkgconfig.generate(libraries: pokipoki_lib,
                   version: '1.0',
                   name: 'libpokipoki',
                   filebase: 'pokipoki',
                   description: 'Object persistence framework')