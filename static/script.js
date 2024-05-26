document.getElementById('orderForm').addEventListener('submit', async function(event) {
	event.preventDefault();
	const orderId = document.getElementById('orderId').value;
	const response = await fetch(`/order/${orderId}`);
	const data = await response.json();
	console.log(data)
	if (data["status"] === "OK") {
		document.getElementById("resp").innerText = JSON.stringify(data["orderItem"], null, 2);
	} else {
		document.getElementById("resp").innerText = data["error"];
	}
});
