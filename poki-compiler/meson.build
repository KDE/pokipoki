go = find_program('go')
sh = find_program('build.sh')
custom_target = custom_target('pokic',
                              build_by_default: true,
                              build_always_stale: true,
                              input: 'go.mod',
                              output: 'pokic',
                              command: [sh, '@INPUT@', '@OUTPUT@'])