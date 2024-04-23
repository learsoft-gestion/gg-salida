import { prefijoURL } from './variables.js';

$(document).ready(function () {
    $('[type=checkbox]').prop('checked', true);

    $('#menu').load('/static/menu.html', function() {
        $('#titulo').append('Archivos');
    });

    flatpickr('.flatpickr', {
        dateFormat: 'd-m-Y',
        locale: 'es',
    });

    $('#btnBuscar').click(function () {
        buscar();
    });

    $('th').click(function () {
        var th = $(this);
        var columnIndex = th.index();
        var order = (th.hasClass('asc')) ? 'desc' : 'asc';

        // Eliminar clases 'asc' y 'desc' de todos los encabezados, y eliminar flechas
        $('th').removeClass('asc desc');
        $('th').find('i.material-icons').remove();

        th.addClass(order);
        th.append('<i class="material-icons expand-icon">' + (order === 'asc' ? 'expand_more' : 'expand_less') + '</i>');

        ordenarTabla(columnIndex, order);
    });

    function buscar() {
        var filters = getFilters();
        $.ajax({
            url: prefijoURL + '/migrador/procesos',
            method: 'GET',
            dataType: 'json',
            data: filters,
            success: function (data) {
                if (data && data.length > 0) {
                    $("#tablaDatos").show();
                    llenarTabla(data);
                } else {
                    console.log('No se recibieron datos del servidor.');
                }
            },
            error: function (error) {
                console.error('Error en la búsqueda:', error);
            }
        })
    }

    function getFilters() {
        json = {
            "idDesde": $('#idDesde').val(),
            "idHasta": $('#idHasta').val(),
            "fechaDesde": $('#fechaInicio').val(),
            "fechaHasta": $('#fechaFin').val(),
            "procesado": "",
            "descripcion": ""
        };

        var procesadoTrue = $('#procesadoTrue').is(':checked');
        var procesadoFalse = $('#procesadoFalse').is(':checked');

        if ( !(procesadoTrue && procesadoFalse) ) {
            if (procesadoTrue || procesadoFalse) {
                json.procesado = procesadoTrue ? true : false;
            }
        }

        var correctoTrue = $('#correctoTrue').is(':checked');
        var correctoFalse = $('#correctoFalse').is(':checked');

        if ( !(correctoTrue && correctoFalse) ) {
            if (correctoTrue || correctoFalse) {
                json.descripcion = correctoTrue ? true : false;
            }
        }

        return json;
    }

    function llenarTabla(data) {
        var tbody = $('table tbody');
        tbody.empty();

        $.each(data, function (index, item) {
            var fechaFormat = "";
            if (item.FechaProcesado.Valid) {
                var fechaProcesado = moment(item.FechaProcesado.Time);
                fechaFormat = fechaProcesado.format("DD/MM/YYYY HH:mm:ss");
            }
            var icon = item.Descripcion.String == "Procesado correctamente" ? "task_alt" : "cancel";
            var iconClass = item.Descripcion.String == "Procesado correctamente" ? "check-icon" : "cancel-icon";
            var row = $('<tr>');
            row.append('<td>' + item.ID + '</td>');
            row.append('<td>' + (item.Procesado ? 'Sí' : 'No') + '</td>');
            row.append('<td>' + fechaFormat + '</td>');
            row.append('<td>' + (item.RutaEntrada.String != "" ? item.RutaEntrada.String : item.ArchivoEntrada) + '</td>');
            row.append('<td>' + (item.RutaFinal.String != "" ? item.RutaFinal.String : item.ArchivoFinal) + '</td>');
            var iconCell = $('<td>');
            iconCell.html('<i class="material-icons ' + iconClass + '" title="' + item.Descripcion.String + '">' + icon + '</i>');
            row.append(iconCell);

            tbody.append(row);
        });
    }

    function ordenarTabla(columnIndex, order) {
        var tbody = $('#tablaDatos tbody');
        var rows = tbody.find('tr').toArray();

        rows.sort(function (a, b) {
            var aValue = $(a).find('td').eq(columnIndex).text();
            var bValue = $(b).find('td').eq(columnIndex).text();

            if (columnIndex === 0) {
                return (order === 'asc') ? parseFloat(aValue) - parseFloat(bValue) : parseFloat(bValue) - parseFloat(aValue);
            } else if (columnIndex === 2) {
                var momentA = moment(aValue, 'DD/MM/YYYY HH:mm:ss');
                var momentB = moment(bValue, 'DD/MM/YYYY HH:mm:ss');

                return (order === 'asc') ? momentA - momentB : momentB - momentA;
            } else {
                return (order === 'asc') ? aValue.localeCompare(bValue) : bValue.localeCompare(aValue);
            }
        });

        tbody.empty();
        $.each(rows, function (index, row) {
            tbody.append(row);
        });
    }
});
