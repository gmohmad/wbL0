document.getElementById('orderForm').addEventListener('submit', async function(event) {
	event.preventDefault();
	const orderId = document.getElementById('orderId').value;
	const response = await fetch(`/order/${orderId}`);
	const data = await response.json();
	document.getElementById("status").innerText = data["status"]
	if (data["status"] === "OK") {
		document.getElementById("orderItem").innerText = JSON.stringify(data["orderItem"], null, 2);
		document.getElementById("errorMsg").innerText = "";
	} else {
		document.getElementById("orderItem").innerText = "";
		document.getElementById("errorMsg").innerText = data["error"];
	}
});
