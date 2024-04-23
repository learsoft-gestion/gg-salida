import { prefijoURL } from './variables.js';

$('#menu').load('/static/menu.html', function() {
    $('#titulo').append('Modelos');
});

// Chequear los checkbox
$('[type=checkbox]').prop('checked', true);

// Select de convenio
$.ajax({
    url: prefijoURL + '/convenios',
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
        url: prefijoURL + `/empresas/${convId}`,
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
        url: prefijoURL + `/conceptos/${convId}/${empId}`,
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
    var json = {
        convenio: $("#conv").val(),
        empresa: $("#emp").val(),
        concepto: $("#conc").val(),
        tipo: $("#tipo").val(),
        jurisdiccion: $("#jurisdiccion").val(),
    }

    var habilitadoTrue = $('#habilitadoTrue').is(':checked');
    var habilitadoFalse = $('#habilitadoFalse').is(':checked');

    if (!(habilitadoTrue && habilitadoFalse)) {
        if (habilitadoTrue || habilitadoFalse) {
            json.vigente = habilitadoTrue ? true : false;
        }
    }
    // Llamada al servidor para mostrar tabla
    $.ajax({
        url: prefijoURL + `/modelos`,
        method: 'GET',
        dataType: 'json',
        data: json,
        success: function (data) {
            if (data && data.length > 0) {
                llenarTabla(data);
            } else {
                $("#tablaDatos").hide();
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

// Función para llenar la tabla
llenarTabla = function (data) {
    $("#tablaDatos").show();
    var tbody = $('table tbody');
    tbody.empty();
    $.each(data, function (index, item) {
        var row = $(`<tr class="accordion-toggle">`);
        row.append(`<td>${item.Convenio}</td>`);
        row.append(`<td>${item.EmpReducido}</td>`);
        row.append(`<td>${item.Concepto}</td>`);
        row.append(`<td>${item.Tipo}</td>`);
        row.append(`<td>${item.Nombre}</td>`);
        row.append(`<td><button class="btn btn-default btn-sm openOculto" data-target=".${item.Id_modelo}"><span class="material-symbols-outlined">arrow_drop_down</span></button></td>`);
        row.append(`<td>${item.Archivo_control}</td>`);
        row.append(`<td>${item.Archivo_modelo}</td>`);
        row.append(`<td>${item.Archivo_nomina}</td>`);
        row.append(`<td><div class="form-check form-switch"><input type="checkbox" value="${item.Id_modelo}" class="form-check-input" ${item.Vigente === "true" ? "checked" : ""}></div></td>`);

        tbody.append(row);

        tbody.append(`<tr class="collapse ${item.Id_modelo}" style="background-color: lightyellow;"><td colspan=12 style="text-align: left;padding-left: 4rem;"><strong>Filtro Convenio:</strong> ${item.Filtro_convenio}</td>`);
        tbody.append(`<tr class="collapse ${item.Id_modelo}" style="background-color: lightyellow;"><td colspan=12 style="text-align: left;padding-left: 4rem;"><strong>Filtro Having:</strong> ${item.Filtro_having}</td>`);
        tbody.append(`<tr class="collapse ${item.Id_modelo}" style="background-color: lightyellow;"><td colspan=12 style="text-align: left;padding-left: 4rem;"><strong>Filtro Personas:</strong> ${item.Filtro_personas}</td>`);
        tbody.append(`<tr class="collapse ${item.Id_modelo}" style="background-color: lightyellow;"><td colspan=12 style="text-align: left;padding-left: 4rem;"><strong>Filtro Recibos:</strong> ${item.Filtro_recibos}</td>`);
    });

    $("tr.accordion-toggle .openOculto").on('click', function () {
        id = $(this).attr("data-target");
        $(id).toggleClass("collapse");
    });

    $('.form-check-input').change(function() {
        var json = {
            id: Number($(this).val()),
            vigente: $(this).is(':checked')
        }
        $.ajax({
            url: prefijoURL + `/modelos`,
            method: 'PATCH',
            dataType: 'json',
            data: JSON.stringify(json),
            success: function (data) {
                Swal.fire({
                    title: "Éxito!",
                    text: data.mensaje,
                    icon: "success"
                });
            },
            error: function (error) {
                if (error.status === 200) {
                    Swal.fire({
                        title: "Éxito!",
                        text: error.statusText,
                        icon: "success"
                    });
                } else {
                    Swal.fire({
                        title: "Ocurrió un error",
                        text: error.mensaje,
                        icon: "error"
                    });
                    console.error('Error en la solicitud:', error);
                }
            }
        })
    });
}