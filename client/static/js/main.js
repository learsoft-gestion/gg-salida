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
    // Llamada al servidor
    $.ajax({
        url: `/procesos`,
        method: 'GET',
        dataType: 'json',
        data: json,
        success: function (data) {
            if (data && data.length > 0) {
                llenarTabla(data);
            } else {
                console.log('No se recibieron datos del servidor.');
            }
        },
        error: function (error) {
            console.error('Error en la búsqueda:', error);
        }
    });
});

// Función para llenar la tabla
llenarTabla = function (data) {
    $("#tablaDatos").show();
    var tbody = $('table tbody');
    tbody.empty();

    $.each(data, function (index, item) {
        var row = $('<tr>');
        row.append('<td>' + item.Empresa + '</td>');
        row.append('<td>' + item.Concepto + '</td>');
        row.append('<td>' + item.Tipo + '</td>');
        row.append('<td>' + item.Nombre + '</td>');
        row.append('<td>' + item.Nombre_salida.String + '</td>');
        row.append('<td>' + item.Version + '</td>');
        row.append('<td>' + item.Ultima_ejecucion + '</td>');
        row.append('<td>' + generarBoton(item) + '</td>');

        tbody.append(row);
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
                    alert(data.mensaje);
                    $("#btnBuscar").trigger("click");
                } else {
                    console.log('No se recibieron datos del servidor.');
                }
            },
            error: function (error) {
                $('#loadingOverlay').hide();
                console.error('Error en la solicitud:', error);
            }
        });
    });
}

generarBoton = function (item) {
    if (item.Boton === "lanzar") {
        return `<button type="button" class="btn btn-success btn-sm lanzar" value="${item.Id}" title="Lanzar"><i class="material-icons">play_arrow</i></button>`;
    } else if (item.Boton === "relanzar") {
        return `<button type="button" class="btn btn-primary btn-sm lanzar" value="${item.Id}" title="Relanzar"><i class="material-icons">refresh</i></button>`;
    }
    return '<button type="button" class="btn btn-transparent" style="width: 38px; height: 38px;" disabled></button>';
}