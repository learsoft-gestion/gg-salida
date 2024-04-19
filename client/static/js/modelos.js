$('#menu').load('/static/menu.html');

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
    var json = {
        convenio: $("#conv").val(),
        empresa: $("#emp").val(),
        concepto: $("#conc").val(),
        tipo: $("#tipo").val(),
        jurisdiccion: $("#jurisdiccion").val(),
    }
    // Llamada al servidor para mostrar tabla
    $.ajax({
        url: `/modelos`,
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
        var row = $(`<tr>`);
        row.append(`<td>${item.Convenio}</td>`);
        row.append(`<td>${item.Empresa}</td>`);
        row.append(`<td>${item.Concepto}</td>`);
        row.append(`<td>${item.Tipo}</td>`);
        row.append(`<td>${item.Nombre}</td>`);
        row.append(`<td><textarea>${item.Filtro_personas}</textarea></td>`);
        row.append(`<td><textarea>${item.Filtro_recibos}</textarea></td>`);
        row.append(`<td><textarea>${item.Filtro_having}</textarea></td>`);
        row.append(`<td><div class="form-check form-switch"><input type="checkbox" value="${item.Id_modelo}" class="form-check-input" ${item.Vigente === "true" ? "checked" : ""}></div></td>`);

        tbody.append(row);
    });

    $('.form-check-input').change(function() {
        // console.log($(this).val());
        // console.log($(this).is(':checked'));
        var json = {
            id: Number($(this).val()),
            vigente: $(this).is(':checked')
        }
        $.ajax({
            url: `/modelos`,
            method: 'PATCH',
            dataType: 'json',
            data: json,
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
        })
    });
}