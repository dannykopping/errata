const signUpButton = document.getElementById('signUp');
const signInButton = document.getElementById('signIn');
const signInSubmit = document.getElementById('signInSubmit');
const container = document.getElementById('container');

signUpButton.addEventListener('click', () => {
	container.classList.add("right-panel-active");
});

signInButton.addEventListener('click', () => {
	container.classList.remove("right-panel-active");
});

signInSubmit.addEventListener('click', () => {
	alert("oops")
});