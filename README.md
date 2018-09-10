UPDATED


This is a work in progress. Things may change by the hour. 
Use at you own risk.

Build image for CrimeData webservice


docker build --no-cache -t crimedata

Run crimedata image


docker run --rm -ti -p 3000:3000 crimedata 

To access the webservice use a command simular to this:


192.168.7.190:3000/crimebook


(subsititue the appropriate ip address)


This will list all crimedata records


To display single records perform:


192.168.7.190:3000/crimebook/n


where n can range between 0-9


I am currently using a smaller crimedata.csv than was issued. My old laptop couldn't handle the larger file size. Currently investigating.


You can also issue PUTs and DELETEs when using a REST client. I'm currently using the POSTMAN REST Client in the Chrome browser.


FILTERING


I've added in some minimal filtering. I've allowed filtering of the data by only one field at a time. Only the following fields can be filtered on: IncidentNumber, OffenseCode, OffenseCodeGroup and District.


EXAMPLES


192.168.7.190:3000/crimebook?District=B2


192.168.7.190:3000/crimebook?OffenseCodeGroup=Larceny


192.168.7.190:3000/crimebook?OffenseCode=619


192.168.7.190:3000/crimebook?IncidentNumber=I182054378


