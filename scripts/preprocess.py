import csv
import os

if __name__ == '__main__':
    # open the file
    cwd = os.getcwd()
    weather = list()

    with open(cwd + '/../src/data/dago.csv') as f:
        # read the file
        reader = csv.reader(f)
        # iterate over the rows
        for row in reader:
            # print the row
            weather.append(row[1])

    # categorize the weather
    categorizer = {
		"Heavy rain":                   "Rainy",
		"Moderate rain":                "Rainy",
		"Light rain":                   "Rainy",
		"Thunderstorm with heavy rain": "Rainy",
		"Haze":                         "Clear",
		"Broken clouds":                "Clear",
		"Scattered clouds":             "Clear",
		"Few clouds":                   "Clear",
		"Overcast clouds":              "Clear",
		"Fog":                          "Clear",
		"Clear Sky":                    "Clear",
    }

    # update the weather
    for i in range(len(weather)):
        weather[i] = categorizer[weather[i]]
    
    from_to = dict() 
    for i in range(1, len(weather)):
        if (weather[i-1], weather[i]) in from_to:
            from_to[(weather[i-1], weather[i])] += 1
        else:
            from_to[(weather[i-1], weather[i])] = 1
    
    # OK! Now we have the data, let's print it
    print(from_to)


    