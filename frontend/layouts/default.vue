<!--
  Файл: layouts/default.vue
  Назначение: основной макет для публичных страниц сайта; подключает общую шапку (AppHeader) и рендерит содержимое страницы через <slot />.
-->
<template>
  <div>
    <AppHeader />
    <slot />
  </div>
</template>
