# dependencies, scripts
- use just command pnpm to execute and install dependencies

# structure
- the directory must be follow the screaming architecture
- the testing must be creating in the directory test in the root directory of the project
- the files must be a spec.ts
- the directory 'test', must be allow the same structure of the directory 'src'
- the domain features must be were in the feature src/features 
- global components and funcionality must be in src/shared
- every features can have pages, components, services, store and directory necessary

# components
- the UI must be created using vuetify, tailwinds 3
- the icons must be using Material Design Icons (mdi)
- every component must be have a test
- if the component need a composable the composable must be have a test
- only create testing with funcionality core 
- the stores is created using pinia

# testing 
- the test must be write using vitetest
- the test must be evalute cases positive and negative
- the test must be have more than one value of case test and maximum 5 values test
- the test must be the most atomic possible
