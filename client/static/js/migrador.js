import { prefijoURL } from './variables.js';

$(document).ready(function () {
    $('[type=checkbox]').prop('checked', true);

    $('#menu').load('/static/menu.html', function () {
        $('#titulo').append('Archivos');
    });

    // Inicializo calendario de Procesado D y H
    flatpickr('.flatpickr', {
        dateFormat: 'd-m-Y',
        locale: 'es',
    });

    // Select de convenio
    $.ajax({
        url: prefijoURL + '/migrador/convenios',
        method: 'GET',
        dataType: 'json',
        success: function (data) {
            if (data && data.length > 0) {
                data.forEach(convenio => {
                    const option = document.createElement("option");
                    option.textContent = convenio;
                    $("#conv").append(option);
                });
            } else {
                console.log('No se recibieron datos del servidor.');
            }
        },
        error: function (error) {
            console.error('Error en la búsqueda:', error);
        }
    });

    // Select de empresa
    $.ajax({
        url: prefijoURL + `/migrador/empresas`,
        method: 'GET',
        dataType: 'json',
        success: function (data) {
            $("#emp").empty();
            var selOption = document.createElement("option");
            selOption.value = '-1';
            selOption.textContent = 'Todas';
            $("#emp").append(selOption);
            if (data && data.length > 0) {
                data.forEach(empresa => {
                    const option = document.createElement("option");
                    option.textContent = empresa;
                    $("#emp").append(option);
                });
            } else {
                console.log('No se recibieron datos del servidor.');
            }
        },
        error: function (error) {
            console.error('Error en la búsqueda:', error);
        }
    });

    // Select de periodo
    $.ajax({
        url: prefijoURL + '/migrador/periodos',
        method: 'GET',
        dataType: 'json',
        success: function (data) {
            if (data && data.length > 0) {
                data.forEach(periodo => {
                    const option = document.createElement("option");
                    option.textContent = periodo;
                    $("#periodo").append(option);
                });
            } else {
                console.log('No se recibieron datos del servidor.');
            }
        },
        error: function (error) {
            console.error('Error en la búsqueda:', error);
        }
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
        var json = {
            "idDesde": $('#idDesde').val(),
            "idHasta": $('#idHasta').val(),
            "fechaDesde": $('#fechaInicio').val(),
            "fechaHasta": $('#fechaFin').val(),
            "procesado": "",
            "descripcion": ""
        };

        if ($('#emp').val() != -1) json.empresa = $('#emp option:selected').text();
        if ($('#conv').val() != -1) json.convenio = $('#conv option:selected').text();
        if ($('#periodo').val() != -1) json.periodo = $('#periodo option:selected').text();

        var procesadoTrue = $('#procesadoTrue').is(':checked');
        var procesadoFalse = $('#procesadoFalse').is(':checked');

        if (!(procesadoTrue && procesadoFalse)) {
            if (procesadoTrue || procesadoFalse) {
                json.procesado = procesadoTrue ? true : false;
            }
        }

        var correctoTrue = $('#correctoTrue').is(':checked');
        var correctoFalse = $('#correctoFalse').is(':checked');

        if (!(correctoTrue && correctoFalse)) {
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
            row.append('<td>' + item.Empresa + '</td>');
            row.append('<td>' + item.Periodo + '</td>');
            row.append('<td>' + item.Convenio + '</td>');
            row.append('<td>' + (item.Procesado ? 'Sí' : 'No') + '</td>');
            row.append('<td>' + fechaFormat + '</td>');
            row.append('<td>' + (item.RutaEntrada.String != "" ? item.RutaEntrada.String : item.ArchivoEntrada) + '</td>');
            row.append('<td>' + (item.RutaFinal.String != "" ? item.RutaFinal.String : item.ArchivoFinal) + '</td>');
            var iconCell = $('<td>');
            iconCell.html('<i class="material-icons ' + iconClass + '" title="' + item.Descripcion.String + '">' + icon + '</i>');
            if (icon === "cancel") {
                var button = `<button type="button" class="btn procesar-btn" value="${item.ID}" title="Procesar"><i class="material-icons">play_arrow</i></button>`;
                iconCell.append(button);
            }
            row.append(iconCell);

            tbody.append(row);
        });

        $('.procesar-btn').click(function() {
            $.ajax({
                url: prefijoURL + `/archivos/${$(this).val()}`,
                method: 'PATCH',
                dataType: 'json',
                success: function (data) {
                    if (data) {
                        Swal.fire({
                            title: "Éxito!",
                            text: data.mensaje,
                            icon: "success"
                        });
                        $("#btnBuscar").trigger("click");
                    } else {
                        console.log('No se pudo procesar el archivo.');
                    }
                },
                error: function (error) {
                    Swal.fire({
                        title: "Ocurrió un error",
                        text: error.mensaje,
                        icon: "error"
                    });
                    console.error('Error en la búsqueda:', error);
                }
            });
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
