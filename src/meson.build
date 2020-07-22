src = [
  'Database.cpp',
  'Database+SQL.cpp',
]

headers = [
  'headers/Database.h',
]

private_headers = [
  'headers/Database+Private.h',
]

moc_files = qt5.preprocess(
  moc_headers: headers,
)

private_moc_files = qt5.preprocess(
    moc_headers: private_headers,
)

install_headers(headers)

fimbeb_lib = shared_library('fimbeb',
                             src,
                             headers,
                             moc_files,
                             private_moc_files,
                             dependencies: [
                                 qt5_deps,
                             ],
                             include_directories: fimbeb_headers,
                             install: true)

pkgconfig = import('pkgconfig')
pkgconfig.generate(libraries: fimbeb_lib,
                   version: '1.0',
                   name: 'libfimbeb',
                   filebase: 'fimbeb',
                   description: 'Object persistence framework')