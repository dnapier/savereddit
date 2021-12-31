# Install packages listed in `package.json` to `node_modules` directory
```bash
npm install
```
# Development
```bash
cd dashboard
npm run start
```
# Production
```bash
cd dashboard
npm run build  # Saves production build to ./dashboard/build directory
cd ..
go run . # Runs Go-gin server that hosts above directory
```