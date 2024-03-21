// Chequear los checkbox
$('[type=checkbox]').prop('checked', true);

// Llenar fecha hasta = fecha desde
$(document).ready(function () {
    // Filtro de fecha/período
    $('#filtroFechaInicio').on('change', function () {
        if ($("#filtroFechaFin").val() === '') {
            $("#filtroFechaFin").val($(this).val());
        }
    });

    $("#filtroFechaInicio, #filtroFechaFin").datepicker({
        autoclose: true,
        minViewMode: 1,
        format: 'mm/yyyy',
        language: "es"
    });

});

// Select de convenio
$.ajax({
    url: '/convenios',
    method: 'GET',
    dataType: 'json',
    success: function (data) {
        if (data && data.length > 0) {
            data.forEach(convenio => {
                const option = document.createElement("option");
                option.value = convenio.id;
                option.textContent = convenio.nombre;
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
$("#conv").change(function () {
    var convId = $("#conv").val();

    $.ajax({
        url: `/empresas/${convId}`,
        method: 'GET',
        dataType: 'json',
        success: function (data) {
            $("#emp").empty();
            var selOption = document.createElement("option");
            selOption.value = '';
            selOption.textContent = 'Todas';
            $("#emp").append(selOption);
            if (data && data.length > 0) {
                data.forEach(empresa => {
                    const option = document.createElement("option");
                    option.value = empresa.id;
                    option.textContent = empresa.nombre;
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
});

// Select de concepto y tipo
$("#emp").change(function () {
    var convId = $("#conv").val();
    var empId = $("#emp").val();

    $.ajax({
        url: `/conceptos/${convId}/${empId}`,
        method: 'GET',
        dataType: 'json',
        success: function (data) {
            $("#conc").empty();
            $("#tipo").empty();
            var selOption = document.createElement("option");
            selOption.value = '';
            selOption.textContent = 'Todos';
            $("#conc").append(selOption);
            var selOption2 = document.createElement("option");
            selOption2.value = '';
            selOption2.textContent = 'Todos';
            $("#tipo").append(selOption2);
            if (data) {
                data.Conceptos.forEach(concepto => {
                    const option = document.createElement("option");
                    option.value = concepto.Id;
                    option.textContent = concepto.Nombre;
                    $("#conc").append(option);
                });
                data.Tipos.forEach(tipo => {
                    const option = document.createElement("option");
                    option.value = tipo.Id;
                    option.textContent = tipo.Nombre;
                    $("#tipo").append(option);
                });
            } else {
                console.log('No se recibieron datos del servidor.');
            }
        },
        error: function (error) {
            console.error('Error en la búsqueda:', error);
        }
    });
});

// Botón buscar
$("#btnBuscar").click(function () {
    // Obtengo todos los valores de los campos en variables
    var conv = $("#conv").val();
    var fechaDesde = $("#filtroFechaInicio").val();
    var fechaHasta = $("#filtroFechaFin").val();
    var json = {
        convenio: conv,
        empresa: $("#emp").val(),
        concepto: $("#conc").val(),
        tipo: $("#tipo").val(),
        jurisdiccion: $("#jurisdiccion").val(),
        fecha1: fechaDesde,
        fecha2: fechaHasta
    }
    var procesadoTrue = $('#procesadoTrue').is(':checked');
    var procesadoFalse = $('#procesadoFalse').is(':checked');

    if (!(procesadoTrue && procesadoFalse)) {
        if (procesadoTrue || procesadoFalse) {
            json.procesado = procesadoTrue ? true : false;
        }
    }
    // Validaciones de campos obligatorios y fechas
    if (!(conv && fechaDesde && fechaHasta)) {
        alert("Faltan completar campos");
        return;
    } else if (!(fechaHasta >= fechaDesde)) {
        alert("La fecha Hasta no puede ser menor a la fecha de inicio");
        return;
    }
    // Muestro mensaje de archivos por generar
    mostrarMensaje(json);
    // Llamada al servidor para mostrar tabla
    $.ajax({
        url: `/procesos`,
        method: 'GET',
        dataType: 'json',
        data: json,
        success: function (data) {
            if (data && data.length > 0) {
                llenarTabla(data);
            } else {
                $("#tablaDatos").hide();
                $("#mensajeFaltantes").hide();
                Swal.fire("No hubo resultados para su búsqueda");
                console.log('No se recibieron datos del servidor.');
            }
        },
        error: function (error) {
            console.error('Error en la búsqueda:', error);
            Swal.fire({
                title: "Ocurrió un error",
                text: error.mensaje,
                icon: "error"
              });
        }
    });
});

mostrarMensaje = function (json) {
    $.ajax({
        url: '/restantes',
        method: 'GET',
        dataType: 'json',
        data: json,
        success: function (data) {
            if (data) {
                $("#mensajeFaltantes").show();
                $("#mensaje").text(data.mensaje);
                if (data.boton) {
                    $("#btnGenerar").show();
                } else {
                    $("#btnGenerar").hide();
                }
            } else {
                console.log('No se recibieron datos del servidor.');
            }
        },
        error: function (error) {
            console.error('Error en la búsqueda:', error);
        }
    });
}

// Función para llenar la tabla
llenarTabla = function (rawData) {
    data = reordenarData(rawData);
    $("#tablaDatos").show();
    var tbody = $('table tbody');
    tbody.empty();
    $.each(data, function (index, item) {
        $.each(item, function (i, proceso) {
            if (i === 0) {
                var row = $(`<tr data-toggle="collapse" class="accordion-toggle">`);
                row.append('<td>' + proceso.Empresa + '</td>');
                row.append('<td>' + proceso.Concepto + '</td>');
                row.append('<td>' + proceso.Tipo + '</td>');
                row.append('<td>' + proceso.Nombre + '</td>');
                row.append(`<td title="${proceso.Nombre_salida.String}"><a href="${proceso.Nombre_salida.String.split("gg-salida")[1]}">${obtenerNombreArchivo(proceso.Nombre_salida.String)}</a></td>`);
                if (proceso.Ultima_version) {
                    row.append(`<td>${proceso.Version}<button class="btn btn-default btn-sm openOculto" data-target="#${proceso.Id}"><span class="material-symbols-outlined">arrow_drop_down</span></button></td>`)
                } else {
                    row.append('<td>' + proceso.Version + '</td>')
                }
                row.append('<td>' + proceso.Ultima_ejecucion + '</td>');
                row.append('<td>' + generarBoton(proceso) + '</td>');

                tbody.append(row);
                if (proceso.Ultima_version) {
                    var subtabla = armarSubtabla(proceso.Id);
                    tbody.append(subtabla);
                }
            } else {
                var subTbody = $(`#tbody-${proceso.Id}`);
                var subRow = $('<tr>');
                subRow.append(`<td title="${proceso.Nombre_salida.String}"><a href="${proceso.Nombre_salida.String.split("gg-salida")[1]}">${obtenerNombreArchivo(proceso.Nombre_salida.String)}</a></td>`);
                subRow.append('<td>' + proceso.Version + '</td>');
                subRow.append('<td>' + proceso.Ultima_ejecucion + '</td>');

                subTbody.append(subRow);
            }
        });
    });

    $("tr.accordion-toggle .openOculto").on('click', function () {
        id = $(this).attr("data-target");
        $(id).toggleClass("collapse");
    });

    // Botones de lanzar y relanzar
    $('.lanzar').click(function () {
        $('#loadingOverlay').show();

        var id = $(this).val();
        var json = {
            Id: $(this).val(),
            Fecha: $("#filtroFechaInicio").val(),
            Fecha2: $("#filtroFechaFin").val()
        }

        $.ajax({
            url: '/send',
            method: 'POST',
            dataType: 'json',
            data: JSON.stringify(json),
            success: function (data) {
                $('#loadingOverlay').hide();
                if (data) {
                    Swal.fire({
                        title: "Éxito!",
                        text: data.mensaje,
                        icon: "success"
                      });
                    $("#btnBuscar").trigger("click");
                } else {
                    console.log('No se recibieron datos del servidor.');
                }
            },
            error: function (error) {
                $('#loadingOverlay').hide();
                Swal.fire({
                    title: "Ocurrió un error",
                    text: error.mensaje,
                    icon: "error"
                  });
                console.error('Error en la solicitud:', error);
            }
        });
    });
}

reordenarData = function (rawData) {
    const data = {};

    rawData.forEach(item => {
        const id = item.Id;
        if (!data[id]) {
            data[id] = [];
        }
        data[id].push(item);
    });

    return data;
}

obtenerNombreArchivo = function (nombre) {
    nombre = nombre.split("\\");
    return nombre[nombre.length - 1];
}

generarBoton = function (item) {
    if (item.Boton === "lanzar") {
        return `<button type="button" class="btn btn-success btn-sm lanzar" value="${item.Id}" title="Lanzar"><i class="material-icons">play_arrow</i></button>`;
    } else if (item.Boton === "relanzar") {
        return `<button type="button" class="btn btn-primary btn-sm lanzar" value="${item.Id}" title="Relanzar"><i class="material-icons">refresh</i></button>`;
    }
    return '<button type="button" class="btn btn-transparent" style="width: 38px; height: 38px;" disabled></button>';
}

armarSubtabla = function (id) {
    var hiddenRow = $('<tr>');
    var td = $('<td colspan="10" class="hiddenRow">');
    var div = $(`<div class="accordian-body collapse" id="${id}">`);
    var table = $('<table class="table">');
    var tbody2 = $(`<tbody id="tbody-${id}">`)
    var tr = $('<tr>');
    tr.append('<th>Nombre de salida</th>');
    tr.append('<th>Versión</th>');
    tr.append('<th>Última ejecución</th>');
    table.append(tr);
    table.append(tbody2);
    div.append(table);
    td.append(div);
    hiddenRow.append(td);

    return hiddenRow;
}

// Botón Generar documentos
$("#btnGenerar").click(function () {
    $('#loadingOverlay').show();

    $.ajax({
        url: '/multiple',
        method: 'POST',
        dataType: 'json',
        success: function (data) {
            $('#loadingOverlay').hide();
            if (data) {
                Swal.fire({
                    title: "Éxito!",
                    text: data.mensaje,
                    icon: "success"
                  });
                $("#btnBuscar").trigger("click");
            } else {
                console.log('No se recibieron datos del servidor.');
            }
        },
        error: function (error) {
            $('#loadingOverlay').hide();
            console.error('Error en la solicitud:', error);
            Swal.fire({
                title: "Ocurrió un error",
                text: error.mensaje,
                icon: "error"
              });
        }
    });
});