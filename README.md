# Medicinal Products Data Fetcher

A simple application for fetching and storing medicinal products data from the Polish e-Health Registry.

## Overview

This application downloads medicinal products data from the Polish e-Health Registry API and stores it locally. It provides a simple interface to view and search through the stored data.

## Features

- Downloads medicinal products data from the Polish e-Health Registry API
- Stores data locally for offline access
- Simple interface for viewing and searching data
- Special API interface for Unitbox integration
- Error handling with fallback to previously downloaded data

## Data Source

The application fetches data from the following endpoint:
```
https://rejestry.ezdrowie.gov.pl/api/rpl/medicinal-products/public-pl-report/6.0.0/overall.xml
```

## Usage Notes

- The application downloads data only once and stores it locally
- No automatic updates are performed
- If a download fails, the application will use the previously downloaded data
- An error will be thrown if no data is available locally

## Technical Details

- The application is designed to be simple and lightweight
- Data is stored locally to minimize API calls
- Error handling ensures graceful degradation when network issues occur

## API Integration

A special API interface is available for Unitbox integration. Please refer to the API documentation for details.