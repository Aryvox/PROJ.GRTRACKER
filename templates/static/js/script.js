const bulle1 = document.querySelector('.bulle1');
const bulle2 = document.querySelector('.bulle2');
const bulle3 = document.querySelector('.bulle3');

document.addEventListener('mousemove', (e) => {

    const mouseX = e.clientX;
    const mouseY = e.clientY;

    bulle1.style.transform = `translate(${mouseX * 0.01}px, ${mouseY * 0.01}px)`;
    bulle2.style.transform = `translate(${mouseX * -0.02}px, ${mouseY * -0.02}px)`;
    bulle3.style.transform = `translate(${mouseX * 0.03}px, ${mouseY * -0.01}px)`;
});
