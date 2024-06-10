/* 
    Aboslutely shit.
    Literally copying and pasting the same functions, horrible code structure.
    But whatever, works.
    Maybe the true clean code are the friends we made along the way.
*/
// *-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-**-*/*-*_*-*

async function resetFirstDropzone() {
    const label = document.querySelector('label[for="dropzone-1-file"]');
    label.classList.remove('border-solid', 'border-blue-500');
    label.classList.add('border-dashed');

    label.querySelector('p').classList.remove('hidden');

    label.querySelector('svg').classList.remove('text-blue-500');
    label.querySelector('svg').classList.add('text-gray-500', 'dark:text-gray-400');

    label.querySelector('p:last-child').classList.remove('hidden');

    label.querySelector('div').classList.remove('p-2');
    label.querySelector('div').classList.add('pt-5', 'pb-6');

    label.style.backgroundImage = '';
}

async function resetSecondDropzone() {
    const label = document.querySelector('label[for="dropzone-2-file"]');
    label.classList.remove('border-solid', 'border-blue-500');
    label.classList.add('border-dashed');

    label.querySelector('p').classList.remove('hidden');

    label.querySelector('svg').classList.remove('text-blue-500');
    label.querySelector('svg').classList.add('text-gray-500', 'dark:text-gray-400');

    label.querySelector('p:last-child').classList.remove('hidden');

    label.querySelector('div').classList.remove('p-2');
    label.querySelector('div').classList.add('pt-5', 'pb-6');

    label.style.backgroundImage = '';
}

async function refreshFirstImage() {
    return new Promise((resolve, reject) => {
        var img = document.getElementById('first-image-display');

        if (fileType === 'tiff' || fileType === 'tif') type = 'jpeg'
        else type = fileType;

        fetch('./api/image-combination?filetype=' + type +
            '&filename=image1' +
            '&_=' + new Date().getTime()) // Prevent caching
            .then(response => response.blob())
            .then(blob => {
                img.src = URL.createObjectURL(blob);
                img.onload = resolve;
                img.onerror = reject;
            })
            .catch(reject);
    });
}

async function refreshSecondImage() {
    return new Promise((resolve, reject) => {
        var img = document.getElementById('second-image-display');

        if (fileType === 'tiff' || fileType === 'tif') type = 'jpeg'
        else type = fileType;

        fetch('./api/image-combination?filetype=' + type +
            '&filename=image2' +
            '&_=' + new Date().getTime()) // Prevent caching
            .then(response => response.blob())
            .then(blob => {
                img.src = URL.createObjectURL(blob);
                img.onload = resolve;
                img.onerror = reject;
            })
            .catch(reject);
    });
}