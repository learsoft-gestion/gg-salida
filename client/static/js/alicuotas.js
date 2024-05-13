var prefijoURL

$('#menu').load('/static/menu.html', function () {
    $('#titulo').append('Alícuotas');
});

fetch("/backend-url")
    .then(res => res.json())
    .then(data => {
        prefijoURL = data.prefijoURL;
        console.log(prefijoURL);

        // Select de Sindicato
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
                        $("#sindicato").append(option);
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
    });

// Llamado a alícuotas
$('#sindicato').change(function () {
    $.ajax({
        url: prefijoURL + `/alicuotas/${$('#sindicato').val()}`,
        method: 'GET',
        dataType: 'json',
        success: function (data) {
            if (data && data.length > 0) {
                mostrarAlicuotas(data);
            } else {
                $('#valores').hide();
                $('#alicuotas').hide();
                Swal.fire("No hubo resultados para su búsqueda");
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

var mostrarAlicuotas = function (data) {
    $.each(data, function (index, alicuota) {
        const group = $(`<div class="form-group">`);
        const control = $(`<div class="form-control">`);
        const nombre = $(`<input type="text" disabled value="${alicuota.Nombre}" id="nombre-${alicuota.IdAlicuota}">`);
        const descripcion = $(`<input type="text" disabled value="${alicuota.Descripcion}" id="descripcion-${alicuota.IdAlicuota}">`);
        const acciones = $(`<div class="actions">`);
        const btnEdit = $(`<button type="button" class="btn btn-sm" id="editAli-${alicuota.IdAlicuota}" title="Editar">`);
        const spanEdit = $(`<span class="material-symbols-outlined">edit</span>`);
        const btnSave = $(`<button type="button" class="btn btn-sm d-none" id="saveAli-${alicuota.IdAlicuota}" title="Guardar">`);
        const spanSave = $(`<span class="material-symbols-outlined">done_outline</span>`);
        const btnDelete = $(`<button type="button" class="btn btn-sm" id="deleteAli-${alicuota.IdAlicuota}" title="Borrar">`);
        const spanDelete = $(`<span class="material-symbols-outlined">delete</span>`);

        btnDelete.append(spanDelete);
        btnSave.append(spanSave);
        btnEdit.append(spanEdit);
        acciones.append(btnEdit);
        acciones.append(btnSave);
        acciones.append(btnDelete);
        control.append(nombre);
        control.append(descripcion);
        control.append(acciones);
        group.append(control);

        $('#alicuotasData').append(group);
        botonEditAli(alicuota.IdAlicuota);
        botonSaveAli(alicuota.IdAlicuota);
        botonDeleteAli(alicuota.IdAlicuota);
    });

    $('#alicuotas').show();
}

$('#btnAddAlicuota').click(function () {
    if ($('#saveAli-').is(':visible')) {
        return;
    }

    const group = $(`<div class="form-group">`);
    const control = $(`<div class="form-control">`);
    const nombre = $(`<input type="text" placeholder="Nombre" id="nombre-">`);
    const descripcion = $(`<input type="text" placeholder="Descripción" id="descripcion-">`);
    const acciones = $(`<div class="actions">`);
    const btnEdit = $(`<button type="button" class="btn btn-sm d-none" id="editAli-" title="Editar">`);
    const spanEdit = $(`<span class="material-symbols-outlined">edit</span>`);
    const btnSave = $(`<button type="button" class="btn btn-sm" id="saveAli-" title="Guardar">`);
    const spanSave = $(`<span class="material-symbols-outlined">done_outline</span>`);
    const btnDelete = $(`<button type="button" class="btn btn-sm d-none" id="deleteAli-" title="Borrar">`);
    const spanDelete = $(`<span class="material-symbols-outlined">delete</span>`);

    btnDelete.append(spanDelete);
    btnSave.append(spanSave);
    btnEdit.append(spanEdit);
    acciones.append(btnEdit);
    acciones.append(btnSave);
    acciones.append(btnDelete);
    control.append(nombre);
    control.append(descripcion);
    control.append(acciones);
    group.append(control);

    $('#alicuotasData').append(group);
});

var botonEditAli = function (id) {
    $(`#editAli-${id}`).click(function () {
        $(`#nombre-${id}`).removeAttr('disabled');
        $(`#descripcion-${id}`).removeAttr('disabled');

        $(`#editAli-${id}`).hide();
        $(`#saveAli-${id}`).removeClass('d-none');

        $.ajax({
            url: prefijoURL + `/valoresAlicuotas/${id}`,
            method: 'GET',
            dataType: 'json',
            success: function (data) {
                if (data && data.length > 0) {
                    mostrarValoresAlicuotas(data);
                } else {
                    $('#valores').hide();
                    Swal.fire("No hubo resultados para su búsqueda");
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
}

var botonSaveAli = function () {

}

var botonDeleteAli = function (id) {
    $(`#deleteAli-${id}`).click(function () {
        Swal.fire({
            title: "¿Quiere borrar la alícuota?",
            text: "No podrá revertir esta acción",
            icon: "warning",
            showCancelButton: true,
            confirmButtonColor: "#3085d6",
            cancelButtonColor: "#d33",
            confirmButtonText: "Sí, borralo",
            cancelButtonText: "Cancelar"
        }).then((result) => {
            if (result.isConfirmed) {
                Swal.fire({
                    title: "¡Borrado!",
                    text: "La alícuota se ha borrado.",
                    icon: "success"
                });
            }
        });
    });
}

var mostrarValoresAlicuotas = function (data) {
    $('#valoresData').empty();

    $.each(data, function (index, valor) {
        const group = $(`<div class="form-group">`);
        const control = $(`<div class="form-control">`);
        const labelPeriodo = $(`<label for="periodo-${valor.IdValoresAlicuota}">Período</label>`);
        const periodo = $(`<input type="text" disabled value="${valor.VigenciaDesde}" id="periodo-${valor.IdValoresAlicuota}">`);
        const labelValor = $(`<br><label for="valor-${valor.IdValoresAlicuota}">Valor</label>`);
        const valorAli = $(`<input type="text" disabled value="${valor.Valor}" id="valor-${valor.IdValoresAlicuota}">`);
        const acciones = $(`<div class="actions">`);
        const btnEdit = $(`<button type="button" class="btn btn-sm" id="editVal-${valor.IdValoresAlicuota}" title="Editar">`);
        const spanEdit = $(`<span class="material-symbols-outlined">edit</span>`);
        const btnSave = $(`<button type="button" class="btn btn-sm d-none" id="saveVal-${valor.IdValoresAlicuota}" title="Guardar">`);
        const spanSave = $(`<span class="material-symbols-outlined">done_outline</span>`);
        const btnDelete = $(`<button type="button" class="btn btn-sm" id="deleteVal-${valor.IdValoresAlicuota}" title="Borrar">`);
        const spanDelete = $(`<span class="material-symbols-outlined">delete</span>`);

        btnDelete.append(spanDelete);
        btnSave.append(spanSave);
        btnEdit.append(spanEdit);
        acciones.append(btnEdit);
        acciones.append(btnSave);
        acciones.append(btnDelete);
        control.append(labelPeriodo);
        control.append(periodo);
        control.append(labelValor);
        control.append(valorAli);
        control.append(acciones);
        group.append(control);

        $('#valoresData').append(group);
        botonEditValor(valor.IdValoresAlicuota);
        botonSaveValor(valor.IdValoresAlicuota);
        botonDeleteValor(valor.IdValoresAlicuota);
    });

    $('#valores').show();
}

var botonEditValor = function (id) {
    $(`#editVal-${id}`).click(function () {
        $(`#periodo-${id}`).removeAttr('disabled');
        $(`#valor-${id}`).removeAttr('disabled');

        $(`#editVal-${id}`).hide();
        $(`#saveVal-${id}`).removeClass('d-none');
    });
}

var botonSaveValor = function (id) {

}

var botonDeleteValor = function (id) {
    $(`#deleteVal-${id}`).click(function () {
        Swal.fire({
            title: "¿Quiere borrar el valor?",
            text: "No podrá revertir esta acción",
            icon: "warning",
            showCancelButton: true,
            confirmButtonColor: "#3085d6",
            cancelButtonColor: "#d33",
            confirmButtonText: "Sí, borralo",
            cancelButtonText: "Cancelar"
        }).then((result) => {
            if (result.isConfirmed) {
                Swal.fire({
                    title: "¡Borrado!",
                    text: "El valor se ha borrado.",
                    icon: "success"
                });
            }
        });
    });
}

$('#btnAddValor').click(function () {
    if ($('#saveVal-').is(':visible')) {
        return;
    }

    const group = $(`<div class="form-group">`);
    const control = $(`<div class="form-control">`);
    const labelPeriodo = $(`<label for="periodo-">Período</label>`);
    const periodo = $(`<input type="text" value="" id="periodo-">`);
    const labelValor = $(`<br><label for="valor-">Valor</label>`);
    const valorAli = $(`<input type="text" value="" id="valor-">`);
    const acciones = $(`<div class="actions">`);
    const btnEdit = $(`<button type="button" class="btn btn-sm d-none" id="editVal-" title="Editar">`);
    const spanEdit = $(`<span class="material-symbols-outlined">edit</span>`);
    const btnSave = $(`<button type="button" class="btn btn-sm" id="saveVal-" title="Guardar">`);
    const spanSave = $(`<span class="material-symbols-outlined">done_outline</span>`);
    const btnDelete = $(`<button type="button" class="btn btn-sm d-none" id="deleteVal-" title="Borrar">`);
    const spanDelete = $(`<span class="material-symbols-outlined">delete</span>`);

    btnDelete.append(spanDelete);
    btnSave.append(spanSave);
    btnEdit.append(spanEdit);
    acciones.append(btnEdit);
    acciones.append(btnSave);
    acciones.append(btnDelete);
    control.append(labelPeriodo);
    control.append(periodo);
    control.append(labelValor);
    control.append(valorAli);
    control.append(acciones);
    group.append(control);

    $('#valoresData').append(group);
});