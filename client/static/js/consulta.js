var prefijoURL

fetch("/backend-url")
    .then(res => res.json())
    .then(data => {
        prefijoURL = data.prefijoURL;
        console.log(prefijoURL);

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
        $.ajax({
            url: prefijoURL + `/empresas`,
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
    })
    .catch(error => {
        console.error("Error al obtener la URL del backend: ", error)
    })

// Chequear los checkbox
$('[type=checkbox]').prop('checked', true);

// Llenar fecha hasta = fecha desde
$(document).ready(function () {
    $('#menu').load('/static/menu.html', function () {
        $('#titulo').append('Informes');
        $('a[href="consulta"]').addClass('active');
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
        filtros.addClass("fixed-top");
    } else {
        filtros.removeClass("fixed-top");
    }
});

// Select de empresas para convenio
$("#conv").change(function () {
    var convId = $("#conv").val();
    var url = convId ? prefijoURL + `/empresas/${convId}` : prefijoURL + "/empresas";

    $.ajax({
        url: url,
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

    var emp = $('#emp').val()
    $.ajax({
        url: prefijoURL + `/jurisdicciones/${convId}${emp ? "/" + emp : "/0"}`,
        method: 'GET',
        dataType: 'json',
        success: function (data) {
            $("#jurisdiccion").empty();
            var selOption = document.createElement("option");
            selOption.value = '';
            selOption.textContent = 'Todas';
            $("#jurisdiccion").append(selOption);
            if (data && data.length > 0) {
                data.forEach(jurisdiccion => {
                    const option = document.createElement("option");
                    option.value = jurisdiccion;
                    option.textContent = jurisdiccion;
                    $("#jurisdiccion").append(option);
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

    $.ajax({
        url: prefijoURL + `/jurisdicciones/${convId ? convId + "/" : "0/"}${empId}`,
        method: 'GET',
        dataType: 'json',
        success: function (data) {
            $("#jurisdiccion").empty();
            var selOption = document.createElement("option");
            selOption.value = '';
            selOption.textContent = 'Todas';
            $("#jurisdiccion").append(selOption);
            if (data && data.length > 0) {
                data.forEach(jurisdiccion => {
                    const option = document.createElement("option");
                    option.value = jurisdiccion;
                    option.textContent = jurisdiccion;
                    $("#jurisdiccion").append(option);
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
            json.consultado = procesadoTrue ? true : false;
        }
    }
    // Validaciones de campos obligatorios y fechas
    if (!(fechaDesde && fechaHasta)) {
        alert("Los campos Desde y Hasta son obligatorios");
        return;
    } else if (!(fechaHasta >= fechaDesde)) {
        alert("La fecha Hasta no puede ser menor a la fecha de inicio");
        return;
    }
    // Llamada al servidor para mostrar tabla
    $.ajax({
        url: prefijoURL + `/consultados`,
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
var llenarTabla = function (rawData) {
    var data = reordenarData(rawData);
    $("#tablaDatos").show();
    var tbody = $('table tbody');
    tbody.empty();
    $.each(data, function (index, item) {
        $.each(item, function (i, proceso) {
            if (i === 0) {
                var row = $(`<tr class="accordion-toggle">`);
                row.append('<td>' + proceso.Convenio + '</td>');
                row.append('<td>' + proceso.Empresa + '</td>');
                row.append('<td>' + proceso.Concepto + '</td>');
                row.append('<td>' + proceso.Tipo + '</td>');
                row.append('<td>' + proceso.Nombre + '</td>');
                proceso.Nombre_control != "-" && proceso.Nombre_control ?
                    row.append(`<td title="${proceso.Nombre_control}"><a href="${obtenerLink(proceso.Nombre_control)}"><span class="material-symbols-outlined">description</span></a></td>`)
                    : row.append('<td>' + proceso.Nombre_control + '</td>');
                proceso.Nombre_nomina != "-" && proceso.Nombre_nomina ?
                    row.append(`<td title="${proceso.Nombre_nomina}"><a href="${obtenerLink(proceso.Nombre_nomina)}"><span class="material-symbols-outlined" style="color: darkorange;">description</span></a></td>`)
                    : row.append('<td>' + proceso.Nombre_nomina + '</td>');
                row.append('<td>' + proceso.Ultima_ejecucion + '</td>');
                row.append('<td>' + generarBoton(proceso.Boton, proceso.Id_modelo, proceso.Id_procesado, "salida") + '</td>');

                tbody.append(row);
            }
        });
    });

    $('table th:nth-child(6), table td:nth-child(6)').css({
        'border-left': '1px solid black',
        'border-right': '1px solid black'
    });
    $('table th:nth-child(8), table td:nth-child(8)').css({
        'border-left': '1px solid black'
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
            url: prefijoURL + '/consulta',
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
                    text: error.responseJSON.mensaje,
                    icon: "error"
                });
                console.error('Error en la solicitud:', error);
                $("#btnBuscar").trigger("click");
            }
        });
    });

    // Botón Eliminar
    eventoBotonDelete();
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
    if (nombre.includes("\\")) {
        nombre = nombre.split("\\");
    } else {
        nombre = nombre.split("/");
    }
    return nombre[nombre.length - 1];
}

var obtenerLink = function (nombre) {
    if (nombre != "" && (nombre.includes("\\api") || nombre.includes("gg-salida"))) {
        return nombre.includes("api") ? prefijoURL + nombre.split("\\api")[1].replace(/\\/g, "/") : prefijoURL + nombre.split("gg-salida")[1].replace(/\\/g, "/");
    } else if (nombre != "") {
        return prefijoURL + nombre;
    }
    return "";
}

var generarBoton = function (boton, id, idProcesado, tipo) {
    if (boton === "lanzar") {
        return `<button type="button" class="btn btn-success btn-sm ${tipo}" value="${id}" title="Consultar"><i class="material-icons">play_arrow</i></button>`;
    }
    return `<button type="button" class="btn btn-primary btn-sm ${tipo}" value="${id}" title="Reconsultar"><i class="material-icons">refresh</i></button>`;
}

var botonDelete = function (boton, id) {
    return boton === "lanzar" ? "" : `<button type="button" class="btn btn-danger btn-sm eliminar" value="${id}" title="Eliminar"><i class="material-icons">delete</i></button>`;
}

var eventoBotonDelete = function () {
    $('.eliminar').click(function () {
        Swal.fire({
            title: "¿Quiere eliminar el proceso?",
            text: "No podrá revertir esta acción",
            icon: "warning",
            showCancelButton: true,
            confirmButtonColor: "#3085d6",
            cancelButtonColor: "#d33",
            confirmButtonText: "Sí, borralo",
            cancelButtonText: "Cancelar"
        }).then((result) => {
            if (result.isConfirmed) {
                $.ajax({
                    url: prefijoURL + '/procesos/' + $(this).val(),
                    method: 'DELETE',
                    dataType: 'json',
                    success: function (data) {
                        if (data) {
                            Swal.fire({
                                title: "Éxito!",
                                text: data,
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
                            text: error.responseJSON.mensaje,
                            icon: "error"
                        });
                        console.error('Error en la solicitud:', error);
                    }
                });
            }
        })
    });
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