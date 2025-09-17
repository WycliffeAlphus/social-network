import React from 'react'

function Events() {
  return (
    <div class="flex h-screen flex-col">
      <header class="flex items-center p-4">
        <button class="mr-2">
          <span class="material-symbols-outlined"> arrow_back_ios </span>
        </button>
        <h1 class="flex-1 text-center text-xl font-bold">Event</h1>
        <div class="w-8"></div>
      </header>
      <main class="flex-1 overflow-y-auto">
        <div class="h-56 w-full bg-cover bg-center" style='background-image: url("https://lh3.googleusercontent.com/aida-public/AB6AXuBQB-gM45XTGTS1XXssoD342ufwIT4K3cEJdwbinOi0ph__SvxZ9BAZwLRh7wu-yfYxclE0hUwMtlTis4GYDCzv2wQj0BombfRbb897BlJ5Oxsfu2lfJb5Wz5oU9cf4ohz_e-D3z3UnKuFFoC_4WtsUd7txMCc4w43Lh09BLjhA5PLT16r2ofEsuNQrvlVzrm3rYz6Tlmzt8q0thpuhI32cPN7H0bwptvCWBrI0VjrqQjADs-uQPvtlBWiAhiEVezl76hHfa2ZXe6I");'></div>
        <div class="p-6">
          <h2 class="text-3xl font-bold">Tech Meetup</h2>
          <p class="mt-2 text-neutral-400">Join us for an evening of networking and discussions on the latest tech trends. Meet fellow enthusiasts and industry experts.</p>
          <div class="mt-6 flex items-center gap-3">
            <span class="material-symbols-outlined text-neutral-400"> calendar_today </span>
            <p class="text-lg">Saturday, July 20</p>
          </div>
          <div class="mt-3 flex items-center gap-3">
            <span class="material-symbols-outlined text-neutral-400"> schedule </span>
            <p class="text-lg">7:00 PM</p>
          </div>
        </div>
      </main>
      <footer class="sticky bottom-0 bg-neutral-900">
        <div class="flex gap-4 p-4">
          <button class="flex-1 rounded-full bg-primary-500 py-3 text-lg font-bold text-white transition-colors hover:bg-opacity-80">Going</button>
          <button class="flex-1 rounded-full bg-neutral-700 py-3 text-lg font-bold text-white transition-colors hover:bg-neutral-600">Not Going</button>
        </div>
        <nav class="flex justify-around border-t border-neutral-700 bg-neutral-800 py-2">
          <a class="flex flex-col items-center gap-1 text-neutral-400" href="#">
            <span class="material-symbols-outlined"> groups </span>
            <span class="text-xs">Groups</span>
          </a>
          <a class="flex flex-col items-center gap-1 text-primary-500" href="#">
            <span class="material-symbols-outlined"> event </span>
            <span class="text-xs">Events</span>
          </a>
          <a class="flex flex-col items-center gap-1 text-neutral-400" href="#">
            <span class="material-symbols-outlined"> person </span>
            <span class="text-xs">Profile</span>
          </a>
        </nav>
      </footer>
    </div>
  )
}

export default Events