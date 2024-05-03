import { prefijoURL } from './variables.js';

// Chequear los checkbox
$('[type=checkbox]').prop('checked', true);

// Llenar fecha hasta = fecha desde
$(document).ready(function () {
    $('#menuContainer').load('/static/menu.html', function() {
        $('#titulo').append('Procesador');
    });
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

var filtros = $("#filtros");
var navbarHeight = $(".navbar").outerHeight(); // Altura de la barra de navegación

$(window).scroll(function () {
    if ($(this).scrollTop() > navbarHeight) {
        // filtros.removeClass('filtros');
        filtros.addClass("fixed-top"); // Agrega la clase para fijar el menú de filtros en la parte superior
    } else {
        // filtros.addClass('filtros');
        filtros.removeClass("fixed-top"); // Quita la clase cuando el usuario se desplaza hacia arriba
    }
});

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
        url: prefijoURL + `/procesos`,
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

var mostrarMensaje = function (json) {
    $.ajax({
        url: prefijoURL + '/restantes',
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
var llenarTabla = function (rawData) {
    var data = reordenarData(rawData);
    $("#tablaDatos").show();
    var tbody = $('table tbody');
    tbody.empty();
    $.each(data, function (index, item) {
        $.each(item, function (i, proceso) {
            if (i === 0) {
                var row = $(`<tr class="accordion-toggle">`);
                row.append('<td>' + proceso.Empresa + '</td>');
                row.append('<td>' + proceso.Concepto + '</td>');
                row.append('<td>' + proceso.Tipo + '</td>');
                row.append('<td>' + proceso.Nombre + '</td>');
                if (proceso.Version > '1') {
                    row.append(`<td>${proceso.Version}<button class="btn btn-default btn-sm openOculto" data-target=".${proceso.Id_modelo}"><span class="material-symbols-outlined">arrow_drop_down</span></button></td>`)
                } else {
                    row.append('<td>' + proceso.Version + '</td>')
                }
                row.append(`<td title="${proceso.Nombre_control}"><a href="${obtenerLink(proceso.Nombre_control)}">${obtenerNombreArchivo(proceso.Nombre_control)}</a></td>`);
                row.append(`<td title="${proceso.Nombre_nomina}"><a href="${obtenerLink(proceso.Nombre_nomina)}">${obtenerNombreArchivo(proceso.Nombre_nomina)}</a></td>`);
                row.append(`<td title="${proceso.Nombre_salida}"><a href="${obtenerLink(proceso.Nombre_salida)}">${obtenerNombreArchivo(proceso.Nombre_salida)}</a></td>`);
                row.append('<td>' + proceso.Ultima_ejecucion + '</td>');
                row.append('<td>' + generarBoton(proceso.Boton, proceso.Id_modelo, proceso.Id_procesado, "salida") + '</td>');

                tbody.append(row);
            } else {
                var subRow = $(`<tr class="collapse ${proceso.Id_modelo}">`);
                subRow.append('<td></td>');
                subRow.append('<td></td>');
                subRow.append('<td></td>');
                subRow.append('<td></td>');
                subRow.append('<td>' + proceso.Version + '</td>');
                subRow.append(`<td title="${proceso.Nombre_control}"><a href="${proceso.Nombre_control.split("gg-salida")[1]}">${obtenerNombreArchivo(proceso.Nombre_control)}</a></td>`);
                subRow.append(`<td title="${proceso.Nombre_nomina}"><a href="${proceso.Nombre_nomina.split("gg-salida")[1]}">${obtenerNombreArchivo(proceso.Nombre_nomina)}</a></td>`);
                subRow.append(`<td title="${proceso.Nombre_salida}"><a href="${proceso.Nombre_salida.split("gg-salida")[1]}">${obtenerNombreArchivo(proceso.Nombre_salida)}</a></td>`);
                subRow.append('<td>' + proceso.Ultima_ejecucion + '</td>');
                subRow.append('<td></td>');

                tbody.append(subRow);
            }
        });
    });

    $('table th:nth-child(6), table td:nth-child(6)').css({
        'border-left': '1px solid black',
        'border-right': '1px solid black'
    });
    $('table th:nth-child(8), table td:nth-child(8)').css({
        'border-left': '1px solid black',
        'border-right': '1px solid black'
    });

    $("tr.accordion-toggle .openOculto").on('click', function () {
        var id = $(this).attr("data-target");
        $(id).toggleClass("collapse");
    });

    // Botones de lanzar y relanzar
    $('.salida').click(function () {
        $('#loadingOverlay').show();

        var json = {
            Id_modelo: Number($(this).val()),
            Fecha: $("#filtroFechaInicio").val(),
            Fecha2: $("#filtroFechaFin").val()
        };

        $.ajax({
            url: prefijoURL + '/send',
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

var reordenarData = function (rawData) {
    const data = {};

    rawData.forEach(item => {
        const id = item.Id_modelo;
        if (!data[id]) {
            data[id] = [];
        }
        data[id].push(item);
    });
    return data;
}

var obtenerNombreArchivo = function (nombre) {
    nombre = nombre.split("\\");
    return nombre[nombre.length - 1];
}

var obtenerLink = function(nombre) {
    return nombre != "" ? nombre.includes("api") ? prefijoURL + nombre.split("api")[1].replace(/\\/g, "/") : prefijoURL + nombre.split("gg-salida")[1].replace(/\\/g, "/") : "";
}

var generarBoton = function (boton, id, idProcesado, tipo) {
    if (boton === "lanzar") {
        return `<button type="button" class="btn btn-success btn-sm ${tipo}" value="${id}" title="Lanzar"><i class="material-icons">play_arrow</i></button>`;
    }
    return `<button type="button" class="btn btn-primary btn-sm ${tipo}" value="${id}" title="Relanzar"><i class="material-icons">refresh</i></button>`;
}

// Botón Generar documentos
$("#btnGenerar").click(function () {
    $('#loadingOverlay').show();

    $.ajax({
        url: prefijoURL + '/multiple',
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

function abrirModal() {
    Swal.fire({
        title: 'Clientes',
        html: `
            <div id="clientesAgregados"></div>
            <input id="cuit" class="swal2-input" placeholder="CUIT">
            <input id="razonSocial" class="swal2-input" placeholder="Razón Social">
            <button type="button" class="btn btn-outline-dark" id="btnBuscarClientes">
              <i class="material-symbols-outlined">search</i>
            </button>
            <div id="resultado"></div>
        `,
        customClass: {
            denyButton: "btn btn-success"
        },
        buttonsStyling: false,
        showConfirmButton: false,
        showDenyButton: true,
        denyButtonText: 'Agregar',
        width: 1000,
        didOpen: () => {
            $('.swal2-deny').hide();
            const denyButton = Swal.getDenyButton();
            $('#btnBuscarClientes').click(function () { buscarCientes(); });
            denyButton.onclick = function () { agregarClientes(); };
        }
    });
}

function buscarCientes() {
    const cuit = $('#cuit').val();
    const razonSocial = $('#razonSocial').val();

    $.ajax({
        url: prefijoURL + '/clientes',
        method: 'GET',
        dataType: 'json',
        data: {
            cuit: cuit,
            cliente: razonSocial
        },
        success: function (data) {
            if (data && data.length > 0) {
                mostrarResultados(data);
            } else {
                $('#resultado').html('<p>No se encontraron resultados</p>');
            }
        },
        error: function (error) {
            console.error('Error en la búsqueda:', error);
        }
    });
}

let clientesAgregados = [];

function agregarClientes() {
    $('.seleccionar-cliente:checked').each(function () {
        const fila = $(this).closest('tr');
        const cuit = fila.find('td:nth-child(2)').text();
        const razonSocial = fila.find('td:nth-child(3)').text();
        const cliente = { cuit: cuit, razonSocial: razonSocial };
        if (!clientesAgregados.some(c => c.cuit === cliente.cuit && c.razonSocial === cliente.razonSocial)) {
            clientesAgregados.push(cliente);
        }
    });
    mostrarClientesAgregados();
}

function mostrarClientesAgregados() {
    let html = '';
    clientesAgregados.forEach(cliente => {
        html += `
            <div class="cliente-agregado">
                <span>${cliente.cuit} / ${cliente.razonSocial}</span>
                <button class="btn btn-danger btn-sm" onclick="quitarCliente('` + cliente.cuit + `')">X</button>
            </div>
        `;
    });
    $('#clientesAgregados').html(html);
}

function quitarCliente(cuit) {
    clientesAgregados = clientesAgregados.filter(cliente => cliente.cuit !== cuit);
    mostrarClientesAgregados();
}

function mostrarResultados(clientes) {
    $('.swal2-deny').show();
    let html = `
        <table class="table">
            <thead>
                <tr>
                    <th><input type="checkbox" id="seleccionarTodo"></th>
                    <th>CUIT</th>
                    <th>Razón Social</th>
                </tr>
            </thead>
            <tbody>
    `;
    clientes.forEach(cliente => {
        html += `
            <tr>
                <td><input type="checkbox" class="seleccionar-cliente" value="${cliente.id}"></td>
                <td>${cliente.cuit}</td>
                <td>${cliente.nombre}</td>
            </tr>
        `;
    });
    html += `
            </tbody>
        </table>
    `;
    $('#resultado').html(html);

    $('#seleccionarTodo').change(function () {
        $('.seleccionar-cliente').prop('checked', $(this).prop('checked'));
    });
}