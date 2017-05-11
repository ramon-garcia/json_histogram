# json_histogram



This program reads JSON data from many files, reads a field from every record, and prints an histogram of frequencies of values of that field.

I use it for classifying events from ELK, when I prefer not to upload local data to ElasticSearch.

To use this program

go build

./json_histogram --field=\<field name\> \<json files\>
