# StockVision

This is the client application that provides the UI for the StockVision project, can search tickers and view their information of ticker and news, recommendations and historical prices.

## Enviroment variables
Create a .env file based on .env.template, fill the values with your own.

``` env
VITE_API_URL=http://localhost:8080
```


## Folder structure
```
├── assets: Assets files
├── features: modules of the application
├── layout: Layout of the application
├── pages: Pages of the application
├── plugins: Plugins to config the application
├── router: Router of the application
└── shared: Shared files that are used in multiple places
    ├──components: Components that are used in multiple places
    ├──composables: Composables that are used in multiple places
    ├──constants: Constants that are used in multiple places
    ├──helpers: utilities that are used in multiple places
    ├──models: Models business logic
    ├──services: services utilities and implementations
    └──store: Pinia store files
```

## Tech stack
- Language: typescript
- Framework: vue 3
- State management: pinia
- CSS framework: tailwindcss 3
- UI components: vuetify

## Installation
You can use `pnpm` or `npm` to execute the commands


```bash
pnpm install
```

```bash
npm install
```

## Run
```bash
pnpm run dev
```

```bash
npm run dev
```

## build
```bash
pnpm run build
```

```bash
npm run build
```

## Commands

**Dev** runs the development server
```bash
pnpm run dev
```

```bash
npm run dev
```

**Build** builds the production version of the app
```bash
pnpm run build
```

```bash
npm run build
```

**Type check** check types of the code to catch errors
```bash
pnpm run type-check
```

```bash
npm run type-check
```

**Preview** visualizes the production build
```bash
pnpm run preview
```

```bash
npm run preview
```


