var personLookup = {};
function personTypeAhead(idField, nameField) {
    $('#' + nameField).typeahead({
        source: function(query, process) {
            $.get({
                url: '/person/json/search?prefix=' + encodeURIComponent(query),
                dataType: 'json',
                success: function(data) {
                    var names = [];
                    $.each(data, function(index, person) {
                        names.push(person.Name);
                        personLookup[person.Name] = person.Id;
                    });
                    process(names);
                },
            });
        },
        updater: function(item) {
          $('#' + idField).val(personLookup[item]);
          return item;
        },
    });
}

var cityLookup = {};
function cityTypeAhead(idField, nameField) {
    $('#' + nameField).typeahead({
        source: function(query, process) {
            $.get({
                url: '/city/json/search?prefix=' + encodeURIComponent(query),
                dataType: 'json',
                success: function(data) {
                    var names = [];
                    $.each(data, function(index, city) {
                        var name = city.Name + ', ' + city.RegionAbbr + ', ' + city.CountryAbbr;
                        names.push(name);
                        cityLookup[name] = city.Id;
                    });
                    process(names);
                },
            });
        },
        updater: function(item) {
          $('#' + idField).val(cityLookup[item]);
          return item;
        },
    });
}
