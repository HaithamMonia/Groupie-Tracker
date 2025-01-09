# Groupie Tracker

Groupie Tracker is a user-friendly website that leverages a provided API to display and interact with data about bands and artists. It offers comprehensive information, including their history, members, concert locations, and dates, enhanced with data visualizations and client-server interactions.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Overview

The Groupie Tracker project involves receiving a given API and manipulating the data contained within it to create a website that displays the information. The API comprises four parts:

1. **Artists**: Contains information about bands and artists, such as their names, images, start years, first album dates, and members.
2. **Locations**: Includes their recent and upcoming concert locations.
3. **Dates**: Lists their recent and upcoming concert dates.
4. **Relation**: Links the data of artists, dates, and locations.

The objective is to build a user-friendly website that displays the band's information through various data visualizations, such as blocks, cards, tables, lists, pages, and graphics. Additionally, the project focuses on creating and visualizing events/actions, particularly implementing client-server interactions where a client action triggers a server response. :contentReference[oaicite:0]{index=0}

## Features

- **Artist Information**: View detailed information about artists, including their history, members, and discography.
- **Concert Details**: Access information about past and upcoming concert dates and locations.
- **Search and Filters**: Utilize a versatile search bar and filtering system to easily find and sort artists.
- **Geolocation**: Visualize concert locations on an interactive map.
- **Data Visualizations**: Experience information presented through various visual formats for better understanding.

## Technologies Used

- **Frontend**:
  - HTML5
  - CSS3
  - JavaScript

- **Backend**:
  - Go (Golang)

- **Data**:
  - JSON-based API

## Installation

Follow these steps to set up the Groupie Tracker project locally:

### Prerequisites

Ensure you have the following installed:

- [Go (Golang)](https://golang.org/dl/) (version 1.16 or higher)
- Git

### Steps

1. **Clone the repository**:
   ```bash
   git clone https://github.com/HaithamMonia/Groupie-Tracker.git
2. **Navigate to the project directory:**
  Groupie-Tracker
3. **Run the application:**
  ```bash
  go run .

