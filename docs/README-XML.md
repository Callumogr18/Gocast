# National Forecast

## Request

```graphql
{
  nationalForecast {
    issued
    today
    outlook
  }
}
```

## Response

```json
{
  "data": {
    "nationalForecast": {
      "region": "National",
      "issued": "2020-06-19T13:39:31+01:00",
      "today": "Through this afternoon and evening, rain over much of Leinster, Ulster, east Munster and east Connacht will gradually clear to the northeast. Brighter weather following from the west and south with some showers. Highest temperatures of 15 to 19 degrees with light to moderate westerly or variable breezes.\r\n",
      "outlook": "Summary: Fresher and breezier this weekend. Changeable with rain at times. \\n\\nSaturday night: Remaining outbreaks of rain in the north and east clearing overnight with clear spells developing. Scattered showers will, however, follow from the Atlantic to affect mainly parts of the west and north overnight. Minimum temperatures of 10 or 11 degrees in moderate to fresh southwest winds, strong near Atlantic coasts...."
    }
  }
}
```

# Provincial Forecast

## Request

```graphql
{
  regionalForecast(region: CONNAUGHT) {
    issued
    tomorrow
    pollen
  }
}
```

## Response

```json
{
  "data": {
    "regionalForecast": {
      "issued": "2020-06-19T11:30:00+01:00",
      "tomorrow": "Any early bright or sunny spells on Saturday will soon give way to increasing cloud from the Atlantic with patchy rain and drizzle developing during the morning. More persistent rain will move in from the southwest during the afternoon, with some heavy bursts in places. It will become windy too, with moderate to fresh and blustery southeasterly winds developing during the morning, becoming strong to near gale in coastal areas. Maximum temperatures of 16 to 18 degrees.\r\n",
      "pollen": "Moderate on Friday and Saturday\r\n"
    }
  }
}
```

# County Forecast

## Request

```graphql
{
  countyForecast(counties: [GALWAY, MAYO]) {
    counties {
      name
      days {
        date
        min_temp
      }
    }
  }
}
```

## Response

```json
{
  "data": {
    "countyForecast": {
      "counties": [
        {
          "name": "GALWAY",
          "days": [
            { "date": "2020-06-19T00:00:00", "min_temp": 9 },
            { "date": "2020-06-20T00:00:00", "min_temp": 11 },
            { "date": "2020-06-21T00:00:00", "min_temp": 10 },
            { "date": "2020-06-22T00:00:00", "min_temp": 8 }
          ]
        },
        {
          "name": "MAYO",
          "days": [
            { "date": "2020-06-19T00:00:00", "min_temp": 9 },
            { "date": "2020-06-20T00:00:00", "min_temp": 11 },
            { "date": "2020-06-21T00:00:00", "min_temp": 10 },
            { "date": "2020-06-22T00:00:00", "min_temp": 8 }
          ]
        }
      ]
    }
  }
}
```
