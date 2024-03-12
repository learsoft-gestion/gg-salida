function mostrarTabla(datos) {
  const tabla = document.createElement('table');
  const thead = document.createElement('thead');
  const tbody = document.createElement('tbody');

  const filaEncabezados = document.createElement('tr');
  const encabezados = ['Nombre', 'Fecha desde', 'Fecha hasta', 'Procesado', 'Cantidad de registros'];
  encabezados.forEach(e => {
    const th = createElement('th');
    th.textContent = e;
    filaEncabezados.appendChild(th);
  });
  thead.appendChild(filaEncabezados);
  tabla.appendChild(thead);

  datos.forEach(registro => {
    const fila = document.createElement('tr');
    const valores = [registro.nombre, registro.fecha_desde, registro.fecha_hasta, registro.procesado, registro.cant_registros];
    valores.forEach(valor => {
      const td = document.createElement('td');
      td.textContent = valor;
      fila.appendChild(td);
    });
    tbody.appendChild(fila);
  });
  tabla.appendChild(tbody);

  const contenedor = document.getElementById('contenedor-tabla');
  contenedor.appendChild(tabla);
}

fetch('/processed')
  .then(res => res.json())
  .then(data => {
    mostrarTabla(data);
  })
  .catch(error => console.error(error));