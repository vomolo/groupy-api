
# Groupie-Trackers-Search-Bar

Groupie Trackers is a web application designed to visualize and interact with data about bands, their concerts, and related information. The project uses data from a provided API to display band information, concert locations, dates, and their relationships in a user-friendly manner.

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Search Functionality](#search-functionality)
- [Installation](#installation)
- [Usage](#usage)
  - [API Integration](#api-integration)
  - [Website Design](#website-design)
  - [Client-Server Interaction](#client-server-interaction)
- [Contributing](#contributing)
- [License](#license)

## Overview
Groupie Trackers pulls data from an API containing:
- **Artists**: Band and artist details including name, image, start year, first album date, and members.
- **Locations**: Concert locations.
- **Dates**: Concert dates.
- **Relation**: Links between artists, locations, and dates.

The project displays this data using various visualizations and allows interaction through client-server communication.

## Features
- **Artist Information**: Display artist profiles with images, names, and details.
- **Concert Information**: Show upcoming and past concerts, including locations and dates.
- **Data Visualization**: Use graphs and charts to visualize concert frequencies, timelines, and other relevant data.
- **Search and Filter**: Dynamic search functionality with typing suggestions, allowing users to search by band name, member name, concert location, first album date, and creation date.
- **Client-Server Interaction**: Features to search and filter concert data based on user input.

## Search Functionality
The application includes a comprehensive search bar that enables users to search and filter through various attributes related to bands and concerts.

- **Case-insensitive Search**: The search bar handles input without considering case sensitivity, making it easier for users to find relevant information.
- **Search Attributes**: The search functionality covers:
  - **Artist/Band Name**
  - **Members** of the band
  - **Concert Locations**
  - **First Album Date**
  - **Creation Date** of the band
- **Typing Suggestions**: As users type in the search bar, suggestions appear, displaying relevant options based on the current input.
  - Each suggestion is labeled with the attribute it matches (e.g., `Freddie Mercury - member` or `Queen - artist/band`).
- **Dynamic Filtering**: Suggestions refine as the user continues typing, making it easier to locate specific data.


## Installation

1. **Clone the repository**:
   ```bash
   git clone https://learn.zone01kisumu.ke/git/antmusumba/groupie-tracker-search-bar
   ```

2. **Navigate to the project directory and install dependencies**:
   ```bash
   cd groupie-trackers-search-bar
   
   ```

## Usage

### API Integration

Fetch data from the provided API endpoints to populate the website:

```javascript
// Example of fetching artist data
fetch('API_URL/artists')
    .then(response => response.json())
    .then(data => {
        // Handle the artist data
    });
```

### Website Design

Design the website to display:
- **Home Page**: Artist profiles using cards or blocks.
- **Concert Info Page**: Tables or lists showing concert locations and dates.
- **Data Visualizations**: Use libraries like Chart.js or D3.js to create graphs and timelines.

### Client-Server Interaction

Implement features that trigger events and communicate with the server. For example, a search feature for concerts:

```javascript
document.querySelector('#searchBtn').addEventListener('click', async () => {
    const artistName = document.querySelector('#artistInput').value;
    const response = await fetch(`API_URL/search?artist=${artistName}`);
    const data = await response.json();
    displayConcerts(data);
});
```

## Contributing

Contributions are welcome! To contribute:
1. Fork the repository.
2. Create a new branch for your feature or fix:
   ```bash
   git checkout -b feature-name
   ```
3. Commit your changes:
   ```bash
   git commit -m 'Add feature or fix'
   ```
4. Push to your branch:
   ```bash
   git push origin feature-name
   ```
5. Open a pull request.

## License

This project is licensed under the MIT License. see the LICENCE file for details
