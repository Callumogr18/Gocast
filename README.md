# Gocast

App that extracts current weather data from `url "https://weather.apis.ie/docs"`, currently provides national, regional and counties weather data.

- XML based data format, the XML is fetched directly from the URL
- XML is queried and unmarshaled into structs which capture the expected response

## Documentation

The docs directory contains...

- README-XML.md which references the query and response format the met.ie API
- FRONTEND-OPTIONS contains plans for possibilty of frontend development, moving from a simple CLI tool
