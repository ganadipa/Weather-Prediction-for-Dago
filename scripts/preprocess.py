import csv
import os


def Weather(cwd: str):
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
    
    # OK! Now we have the data, let's count it
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

def Humidity_Temperature_WindSpeed():
    state = list()

    with open(cwd + '/data/dago-complete.csv') as f:
        # read the file
        reader = csv.reader(f)
        # iterate over the rows
        for row in reader:
            '''
                row[2] is humidity, row[3] is temperature, row[5] is wind speed
            '''
            state.append((row[2], row[3], row[5]))
    
    # categorize the humidity
    for i in range(len(state)):
        if (state[i][0] < 60):
            state[i] = ('Kering', state[i][1], state[i][2])
        elif (state[i][0] < 80):
            state[i] = ('Lembab', state[i][1], state[i][2])
        else:
            state[i] = ('Basah', state[i][1], state[i][2])

    # categorize the temperature
    for i in range(len(state)):
        if (state[i][1] < 20):
            state[i] = (state[i][0], 'Dingin', state[i][2])
        elif (state[i][1] < 30):
            state[i] = (state[i][0], 'Sejuk', state[i][2])
        else:
            state[i] = (state[i][0], 'Panas', state[i][2])
    
    # categorize the wind speed
    for i in range(len(state)):
        if (state[i][2] < 5):
            state[i] = (state[i][0], state[i][1], 'Tenang')
        elif (state[i][2] < 10):
            state[i] = (state[i][0], state[i][1], 'Kencang')
        else:
            state[i] = (state[i][0], state[i][1], 'Sangat Kencang')
    

    from_to = dict()
    for i in range(1, len(state)):
        if (state[i-1], state[i]) in from_to:
            from_to[(state[i-1], state[i])] += 1
        else:
            from_to[(state[i-1], state[i])] = 1
    
    # OK! Now we have the data, let's count it

def generate_state_to_index(from_to: dict):
    num_types_for_each_index = [set() for i in range(len(from_to.keys()))]
    for key, value in from_to.items():
        for i in range(len(key)):
            num_types_for_each_index[i].add(key[i])

def categorize_data(cwd: str):
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

    input_file = os.path.join(cwd, 'data', 'dago.csv')
    output_file = os.path.join(cwd,  'test', 'categorized_dago.csv')

    # read the data and put it in the categorized data csv
    with open(input_file, 'r') as f:
        reader = csv.reader(f)
        with open(output_file, 'w') as g:
            writer = csv.writer(g)
            for row in reader:
                # write row[0] comma categorizer[row[1]]
                writer.writerow([row[0], categorizer[row[1]]])





    

if __name__ == '__main__':
    # open the file
    cwd = os.getcwd()

    # get the weather-only transition matrix
    categorize_data(cwd)

    


    