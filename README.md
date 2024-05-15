# basedQuotes
Install
  Add env files
    Add `.env.local` in `/frontend`
    Add `.env` in base directory

  Launch the NextJS app
    Navigate to `/frontend`
    Install dependencies
      `npm i`
    Launch
      `npm run dev`
  Build the go app
    In the base directory
      `go build`

Usage
`go run ./main.go [Venue] [base]/[quote]`
ex
`go run ./main.go sfox btc/usd`
`go run ./main.go sfox eth/usd`

The modularity hinges on a switch case selecting which implementation of `venue` interface to use. The SFOX package is the only complete implementation. Dropping in a new venue shouldn't require anything outside of adding a new package within `/adapters` and adding it to the switchcase.
