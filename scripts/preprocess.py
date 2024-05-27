import csv
import os

if __name__ == '__main__':
    # open the file
    cwd = os.getcwd()
    weather = list()

    with open(cwd + '/data/dago.csv') as f:
        # read the file
        reader = csv.reader(f)
        # iterate over the rows
        for row in reader:
            # print the row
            weather.append(row[1])

    # categorize the weather
    categorizer = {
        "Clear Sky":                    "Terang",
        "Few clouds":                   "Terang",
        "Scattered clouds":             "Berawan",
        "Broken clouds":                "Berawan",
        "Overcast clouds":              "Mendung",
        "Fog":                          "Mendung",
        "Haze":                         "Mendung",
        "Light rain":                   "Hujan",
        "Moderate rain":                "Hujan",
        "Heavy rain":                   "Hujan",
        "Thunderstorm with heavy rain": "Hujan",
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
    count = [0 for i in range(4)]
    for key, value in from_to.items():
        if (key[0] == 'Terang'):
            count[0] += value
        elif (key[0] == 'Mendung'):
            count[1] += value
        elif (key[0] == 'Hujan'):
            count[2] += value
        else:
            assert key[0] == 'Berawan'
            count[3] += value
    
    probability = dict()
    for key, value in from_to.items():
        if (key[0] == 'Terang'):
            probability[key] = value / count[0]
        elif (key[0] == 'Mendung'):
            probability[key] = value / count[1]
        elif (key[0] == 'Hujan'):
            probability[key] = value / count[2]
        else:
            assert key[0] == 'Berawan'
            probability[key] = value / count[3]
    
    print(probability)

    # State to index
    state_to_index = {
        'Terang': 0,
        'Berawan': 1,
        'Mendung': 2,
        'Hujan': 3
    }

    # Transition matrix
    transition_matrix = [[0 for i in range(4)] for j in range(4)]
    for key, value in probability.items():
        transition_matrix[state_to_index[key[0]]][state_to_index[key[1]]] = value

    # Write to file cwd + '/data/dago-matrix.csv'
    with open(cwd + '/data/dago-matrix.csv', 'w') as f:
        writer = csv.writer(f)
        writer.writerows(transition_matrix)


    