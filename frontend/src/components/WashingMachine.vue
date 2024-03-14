<template>
    <div class="washing-machine-outer-body" :style="washStyle">
        <div class="washing-machine">
            <div :class="{ body: true, 'animate-shake': animate }">
                <div class="panel">
                    <div class="button"></div>
                    <div class="dial"></div>
                </div>
                <div class="door">
                    <div class="window">
                        <div :class="{ drum: true, 'drum-rotate': animate }">
                            <img src="/img/walzen_logo_transparent.png" alt="W" class="letter-W" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
const props = defineProps({
    scale: {
        type: Number,
        default: 1,
    },
    animate: {
        type: Boolean,
        default: false,
    },
});

const washStyle = computed(() => ({
    "--scale-factor": props.scale,
}));
</script>

<style scoped lang="scss">
$base-width: 150px;
$base-height: 150px;
$base-radius: 50%;
$panel-height: 40px;
$panel-margin-top: 20px;
$panel-border-radius: 10px;
$control-size: 20px;
$control-dial-size: 25px;
$base-door-color: #455a64;
$base-panel-color: #37474f;
$control-color: #cfd8dc;
$drum-width: 100px;
$drum-height: 200px;
$drum-color: #303030;
$padding: 20px;
$button-margin: 10px; // Assuming there is a margin around buttons or dials
$w-dimensions: 100px;

@keyframes rotate {
    0% {
        transform: rotate(0deg);
    }

    100% {
        transform: rotate(360deg);
    }
}

@keyframes shake {
    0%,
    100% {
        transform: translate(0, 0) rotate(0);
    }

    25% {
        transform: translate(-0.5px, 0.5px) rotate(0.3deg);
    }

    50% {
        transform: translate(-0.3px, -0.3px) rotate(0.3deg);
    }
}

// Mixin for scaling
@mixin scale($property, $value) {
    #{$property}: calc(#{$value} * var(--scale-factor));
}

.washing-machine {
    display: flex;
    /*
    justify-content: center;
    transform-origin: center;
    */
    align-items: center;
}

.body {
    background-color: #0277bd;
    @include scale(border-radius, 20px);
    @include scale(padding, $padding);
    display: flex;
    flex-direction: column;
    align-items: center;
    position: relative;
}

.animate-shake {
    animation: shake 0.2s linear infinite;
}

.door {
    @include scale(width, $base-width);
    @include scale(height, $base-height);
    background-color: #455a64;
    border-radius: $base-radius;
    display: flex;
    justify-content: center;
    align-items: center;
    position: relative;
    z-index: 1;
}

.window {
    @include scale(width, $base-width - 30px);
    @include scale(height, $base-height - 30px);
    background-color: rgba(255, 255, 255, 0.2);
    border-radius: $base-radius;
    display: flex;
    justify-content: center;
    align-items: center;
    position: relative;
    overflow: hidden;
}

.drum {
    position: relative;
    @include scale(width, $drum-width);
    /* Adjust size as needed */
    @include scale(height, $drum-height);
    background-color: $drum_color;
    /* Drum color */
    border-radius: $base-radius;
    display: flex;
    justify-content: center;
    align-items: center;
    overflow: hidden;
    /* Hide anything that goes outside the drum circle */
}

.drum-rotate {
    animation: rotate 2.5s linear infinite;
}

.letter-W {
    max-width: calc($w-dimensions * var(--scale-factor));
    /* Adjust size as needed */
    max-height: calc($w-dimensions * var(--scale-factor));
    /* Adjust size as needed */
    width: auto;
    height: auto;
}

.panel {
    @include scale(width, 150px);
    @include scale(height, $panel-height);
    background-color: #37474f;
    @include scale(margin-bottom, $panel-margin-top);
    @include scale(border-radius, $panel-border-radius);
    display: flex;
    justify-content: space-around;
    align-items: center;
    padding: 0 calc(20px * var(--scale-factor));
}

.button,
.dial {
    @include scale(width, $control-size);
    @include scale(height, $control-size);
    background-color: #cfd8dc;
    border-radius: 50%;
}

.dial {
    @include scale(width, $control-dial-size);
    @include scale(height, $control-dial-size);
}
</style>
