project('clox', 'c')

sources = ['main.c', 'memory.c', 'debug.c', 'chunk.c', 'value.c']

add_global_arguments('-Wextra', language : 'c')

if get_option('buildtype').startswith('release')
  add_project_arguments('-DDEBUG=0', language : ['c'])
endif
if get_option('buildtype').startswith('debug')
  add_project_arguments('-DDEBUG=1', language : ['c'])
endif


executable('clox', sources)
