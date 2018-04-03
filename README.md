# weather-color

This golang server is fetching informations form OpenWeatherMap (https://openweathermap.org/current)
and calculates the Temp, Humidity and Visibility(pollution) to return a coresponded color according
to the city that has been entered in url.

### Rules for colors : 

- Red = Temp => -20C = 0 / 50C = 255
- Blue = Humidity => 0% = 0 / 50%=123 / +123 if it's raining
- Green = 0 polution = 255 / a lot of pollution 255

pollution according to visibility (PM2.5 and PM10 particles)
see ref. (http://www.epa.vic.gov.au/your-environment/air/air-pollution/visibility-reduction)


### Prerequisites : 

This Server don't use any external Package,
you just need to install and configure GO.

### Test it :

```
Go run server.go (Launch to localhost:9090)
```

### Utility Links :

```
http://localhost:9090/ -> home
http://localhost:9090/Paris -> result
http://localhost:9090/Mykonos -> result
http://localhost:9090/DontExist -> error message
```

UnitTests Are present to test and debug

### Author :
Laurent Loukopoulos