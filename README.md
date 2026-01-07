# MapMyFunds

MapMyFunds is a data visualization tool that displays campaign contribution totals to U.S. presidential candidates at the state level. The application uses publicly available data from the Federal Election Commission via the OpenFEC API ([api.open.fec.gov](https://api.open.fec.gov)) and visualizes aggregated, precompiled contribution totals through interactive geographic maps.

## Methodology

### Data Collection and Storage

MapMyFunds consumes data published by the Federal Election Commission through the OpenFEC API. The `/presidential` endpoint provides state-level and national aggregates of total contributions received by each presidential candidate for a given election cycle. To reduce request volume and improve performance, these records are persisted in a MongoDB document database.

Retrieved data is organized by election year, with contribution totals grouped by state and candidate.

### Serving and Visualization

Contribution totals are exposed through a Go-based REST API and consumed by a Vite–React frontend. The frontend renders the data using D3.js, displaying state-level contribution totals on an interactive map. States are color-coded in an electoral-style format, with aggregate contribution breakdowns shown alongside the visualization.

### Deployment

MapMyFunds is hosted at [https://mapmyfunds.zeeshawnh.com](https://mapmyfunds.zeeshawnh.com) on an AWS EC2 instance and deployed via a GitHub Actions CI/CD pipeline.

## Notes

This project visualizes aggregated contribution totals as reported by the Federal Election Commission. It does not attribute individual donor intent or model campaign finance relationships beyond the published aggregates.
