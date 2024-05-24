var prefijoURL

$('#menu').load('/static/menu.html', function () {
    $('#titulo').append('Personal Interno');
});

fetch("/backend-url")
    .then(res => res.json())
    .then(data => {
        prefijoURL = data.prefijoURL;
        console.log(prefijoURL);
        buscarCuils();
    })
    .catch(error => {
        console.error("Error al obtener la URL del backend: ", error)
    });

var buscarCuils = function () {
    $.ajax({
        url: prefijoURL + '/personalinterno/all',
        method: 'GET',
        dataType: 'json',
        success: function (data) {
            if (data && data.length > 0) {
                mostrarCuils(data);
            }
        },
        error: function (error) {

        }
    });
}

var mostrarCuils = function (data) {
    $('#cuilsData').empty();
    data.forEach(cuil => {
        const group = $(`<div class="form-group" id="${cuil}">`);
        const control = $(`<div class="form-control">`);
        const nombre = $(`<input type="text" disabled value="${cuil}" id="nombre-${cuil}" class="nombre">`);
        const acciones = $(`<div class="actions">`);
        const btnDelete = $(`<button type="button" class="btn btn-sm" id="delete-${cuil}" title="Borrar">`);
        const spanDelete = $(`<span class="material-symbols-outlined">delete</span>`);

        btnDelete.append(spanDelete);
        acciones.append(btnDelete);
        control.append(nombre);
        control.append(acciones);
        group.append(control);

        $('#cuilsData').append(group);
        botonDeleteCuil(cuil);
    });
}

var botonDeleteCuil = function (id) {
    $(`#delete-${id}`).click(function () {
        Swal.fire({
            title: "¿Quiere borrar el cuil?",
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
                    url: prefijoURL + `/personalinterno`,
                    method: 'DELETE',
                    dataType: 'json',
                    data: JSON.stringify({ cuil: String(id) }),
                    success: function (data) {
                        if (data) {
                            buscarCuils();
                            Swal.fire({
                                title: "¡Borrado!",
                                text: "El cuil se ha borrado.",
                                icon: "success"
                            });
                        } else {
                            Swal.fire("No se pudo borrar el cuil");
                        }
                    },
                    error: function (error) {
                        console.error('Error intentando borrar:', error);
                        Swal.fire({
                            title: "Ocurrió un error",
                            text: error.responseText,
                            icon: "error"
                        });
                    }
                });
            }
        });
    });
}

$('#btnAdd').click(function () {
    if ($('#new').is(':visible')) {
        return;
    }

    const group = $(`<div class="form-group" id="new">`);
    const control = $(`<div class="form-control">`);
    const nombre = $(`<input type="text" value="" id="nombre" class="nombre" placeholder="Valor">`);
    const acciones = $(`<div class="actions">`);
    const btnSave = $(`<button type="button" class="btn btn-sm" id="save" title="Guardar">`);
    const spanSave = $(`<span class="material-symbols-outlined">done_outline</span>`);
    const btnCancel = $(`<button type="button" class="btn btn-sm" id="cancel" title="Cancelar">`);
    const spanCancel = $(`<span class="material-symbols-outlined">close</span>`);
    const btnDelete = $(`<button type="button" class="btn btn-sm d-none" id="delete" title="Borrar">`);
    const spanDelete = $(`<span class="material-symbols-outlined">delete</span>`);

    btnDelete.append(spanDelete);
    btnCancel.append(spanCancel);
    btnSave.append(spanSave);
    acciones.append(btnSave);
    acciones.append(btnCancel);
    acciones.append(btnDelete);
    control.append(nombre);
    control.append(acciones);
    group.append(control);

    $('#cuilsData').append(group);
    botonCancel();
    botonSaveCuil();
});

var botonCancel = function () {
    $('#cancel').click(function () {
        $('#new').remove();
    });
}

var botonSaveCuil = function () {
    $('#save').click(function () {
        var cuil = $('#nombre').val();
        if (cuil == "") {
            alert('Debe escribir un valor de cuil.');
            return;
        }
        $.ajax({
            url: prefijoURL + '/personalinterno',
            method: 'POST',
            dataType: 'json',
            data: JSON.stringify({ cuil: cuil }),
            success: function (data) {
                $(`#new`).attr('id', `${cuil}`);
                $(`#nombre`).attr('id', `nombre-${cuil}`);
                $(`#nombre-${cuil}`).prop('disabled', true);
                $(`#save`).attr('id', `save-${cuil}`);
                $(`#delete`).attr('id', `delete-${cuil}`);
                $(`#save-${cuil}`).hide();
                $(`#delete-${cuil}`).removeClass('d-none');
                $(`#cancel`).remove();
                botonDeleteCuil(cuil);

            },
            error: function (error) {
                console.error('Error en la búsqueda:', error);
                Swal.fire({
                    title: "Ocurrió un error",
                    text: error.responseText,
                    icon: "error"
                });
            }
        });
    });
}