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
    $('#menuContainer').load('/static/menu.html', function () {
        $('#titulo').append('Procesador');
        $('#home').addClass('active');
    });
    // Filtro de fecha
    $("#filtroFechaInicio").datepicker({
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
    var fechaHasta = $("#filtroFechaInicio").val();
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
    if (!fechaDesde) {
        alert("Los campos Desde y Hasta son obligatorios");
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
                    // $("#btnGenerar").hide(); // Se deja oculto botón para que no se vea en prod hasta nuevo fix
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
                row.append('<td>' + proceso.Convenio + '</td>');
                row.append('<td>' + proceso.Empresa + '</td>');
                row.append('<td>' + proceso.Concepto + '</td>');
                row.append('<td>' + proceso.Tipo + '</td>');
                row.append('<td>' + proceso.Nombre + '</td>');
                if (item.length > 1) {
                    row.append(`<td>${proceso.Version}<button class="btn btn-default btn-sm openOculto" data-target=".${proceso.Id_modelo}"><span class="material-symbols-outlined">arrow_drop_down</span></button></td>`)
                } else {
                    row.append('<td>' + proceso.Version + '</td>')
                }
                proceso.Nombre_control != "-" && proceso.Nombre_control != "" ?
                    row.append(`<td title="${proceso.Nombre_control}"><a href="${obtenerLink(proceso.Nombre_control)}"><span class="material-symbols-outlined">description</span></a></td>`)
                    : row.append('<td>' + proceso.Nombre_control + '</td>');
                proceso.Nombre_nomina != "-" && proceso.Nombre_nomina != "" ?
                    row.append(`<td title="${proceso.Nombre_nomina}"><a href="${obtenerLink(proceso.Nombre_nomina)}"><span class="material-symbols-outlined" style="color: darkorange;">description</span></a></td>`)
                    : row.append('<td>' + proceso.Nombre_nomina + '</td>');
                proceso.Nombre_salida != "-" && proceso.Nombre_salida != "" ?
                    row.append(`<td title="${proceso.Nombre_salida}"><a href="${obtenerLink(proceso.Nombre_salida)}"><span class="material-symbols-outlined" style="color: green;">description</span></a></td>`)
                    : row.append('<td>' + proceso.Nombre_salida + '</td>');
                row.append('<td>' + proceso.Ultima_ejecucion + '</td>');
                row.append('<td>' + generarBoton(proceso.Boton, proceso.Id_modelo, proceso.Id_procesado, "salida", proceso.Bloqueado) + botonDelete(proceso.Boton, proceso.Id_procesado, proceso.Bloqueado) + '</td>');
                row.append('<td>' + botonBloquear(proceso.Boton, proceso.Id_procesado, proceso.Bloqueado) + '</td>');

                tbody.append(row);
            } else {
                var subRow = $(`<tr class="collapse ${proceso.Id_modelo}">`);
                subRow.append('<td></td>');
                subRow.append('<td></td>');
                subRow.append('<td></td>');
                subRow.append('<td></td>');
                subRow.append('<td></td>');
                subRow.append('<td>' + proceso.Version + '</td>');
                proceso.Nombre_control != "-" && proceso.Nombre_control != "" ?
                    subRow.append(`<td title="${proceso.Nombre_control}"><a href="${obtenerLink(proceso.Nombre_control)}"><span class="material-symbols-outlined">description</span></a></td>`)
                    : subRow.append('<td>' + proceso.Nombre_control + '</td>');
                proceso.Nombre_nomina != "-" && proceso.Nombre_nomina != "" ?
                    subRow.append(`<td title="${proceso.Nombre_nomina}"><a href="${obtenerLink(proceso.Nombre_nomina)}"><span class="material-symbols-outlined" style="color: darkorange;">description</span></a></td>`)
                    : subRow.append('<td>' + proceso.Nombre_nomina + '</td>');
                proceso.Nombre_salida != "-" && proceso.Nombre_salida != "" ?
                    subRow.append(`<td title="${proceso.Nombre_salida}"><a href="${obtenerLink(proceso.Nombre_salida)}"><span class="material-symbols-outlined" style="color: green;">description</span></a></td>`)
                    : subRow.append('<td>' + proceso.Nombre_salida + '</td>');
                subRow.append('<td>' + proceso.Ultima_ejecucion + '</td>');
                subRow.append('<td>' + botonDelete(null, proceso.Id_procesado) + '</td>');
                subRow.append('<td></td>');

                tbody.append(subRow);
            }
        });
    });

    $('table th:nth-child(7), table td:nth-child(7)').css({
        'border-left': '1px solid black',
        'border-right': '1px solid black'
    });
    $('table th:nth-child(9), table td:nth-child(9)').css({
        'border-left': '1px solid black',
        'border-right': '1px solid black'
    });
    $('table th:nth-child(11), table td:nth-child(11)').css({
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

    // Botón Bloquear
    eventoBotonBloquear();
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

var generarBoton = function (boton, id, idProcesado, tipo, bloqueado) {
    if (boton === "lanzar") {
        return `<button type="button" class="btn btn-success btn-sm ${tipo}" value="${id}" title="Lanzar"><i class="material-icons">play_arrow</i></button>`;
    }
    return `<button type="button" class="btn btn-primary btn-sm ${tipo}" value="${id}" title="Relanzar" ${bloqueado ? "disabled" : ""}><i class="material-icons">refresh</i></button>`;
}

var botonDelete = function (boton, id, bloqueado) {
    return boton === "lanzar" ? "" : `<button type="button" class="btn btn-danger btn-sm eliminar" value="${id}" title="Eliminar" ${bloqueado ? "disabled" : ""}><i class="material-icons">delete</i></button>`;
}

var botonBloquear = function (boton, id, bloqueado) {
    return boton === "lanzar" ? "" : `<button type="button" class="btn btn-warning btn-sm bloquear" value="${id}" title="${bloqueado ? 'Desbloquear' : 'Bloquear'}"><i class="material-icons">${bloqueado ? 'lock_open' : 'lock'}</i></button>`;
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

var eventoBotonBloquear = function () {
    $('.bloquear').click(function () {
        var idProceso = $(this).val();
        var json = {
            "Id_procesado": Number(idProceso),
            "Bloquear": $(this).attr('title') === 'Bloquear'
        }
        $.ajax({
            url: prefijoURL + '/send',
            method: 'PATCH',
            data: JSON.stringify(json),
            success: function (data) {
                console.log("Success", data);
                if (data) {
                    $("#btnBuscar").trigger("click");
                }
            },
            error: function (error) {
                console.log("Error: ", error);
                Swal.fire({
                    title: "Ocurrió un error",
                    text: error.responseText,
                    icon: "error"
                });
            }
        });
    })
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