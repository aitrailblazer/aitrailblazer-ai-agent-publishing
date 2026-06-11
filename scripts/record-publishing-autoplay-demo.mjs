import { chromium } from 'playwright';
import fs from 'node:fs/promises';
import path from 'node:path';

const repo = process.cwd();
const timestamp = new Date().toISOString().replace(/[-:]/g, '').replace(/\.\d{3}Z$/, 'Z');
const outDir = path.join(repo, 'artifacts', '03_THREE_MINUTE_PUBLISHING_DEMO', `${timestamp}-publishing-browser-autoplay`);
const rawVideoDir = path.join(outDir, 'playwright-video-raw');
const url = process.env.PUBLISHING_DEMO_URL || 'http://127.0.0.1:48805/demo';
const durationMs = Number(process.env.PUBLISHING_AUTOPLAY_RECORD_MS || '178000');
const landingHoldMs = Number(process.env.PUBLISHING_LANDING_HOLD_MS || '18000');
const width = Number(process.env.PUBLISHING_RECORD_WIDTH || '1280');
const height = Number(process.env.PUBLISHING_RECORD_HEIGHT || '720');

await fs.mkdir(rawVideoDir, { recursive: true });

const browser = await chromium.launch({
  headless: true,
  args: ['--disable-dev-shm-usage', '--hide-scrollbars', '--force-device-scale-factor=1'],
});

let videoPath = '';
try {
  const context = await browser.newContext({
    viewport: { width, height },
    deviceScaleFactor: 1,
    recordVideo: { dir: rawVideoDir, size: { width, height } },
  });
  const page = await context.newPage();
  await page.goto(url, { waitUntil: 'networkidle', timeout: 30000 });
  await page.screenshot({ path: path.join(outDir, '000-start.png'), fullPage: false });
  if (url.endsWith('/demo') || url.endsWith('/demo/')) {
    await page.waitForTimeout(Math.min(landingHoldMs, durationMs));
    await page.locator('a.button[href*="/demo/run?autoplay=1"]').first().click();
    await page.waitForLoadState('domcontentloaded', { timeout: 30000 });
    await page.screenshot({ path: path.join(outDir, '010-demo-start.png'), fullPage: false });
    await page.waitForTimeout(Math.max(0, durationMs - landingHoldMs));
  } else {
    await page.waitForTimeout(durationMs);
  }
  await page.screenshot({ path: path.join(outDir, '999-end.png'), fullPage: false });
  const video = page.video();
  await context.close();
  videoPath = video ? await video.path() : '';
} finally {
  await browser.close();
}

if (!videoPath) {
  throw new Error('Playwright did not produce a video artifact.');
}

const finalVideo = path.join(outDir, 'AITrailblazer_Publishing_Three_Minute_Demo.webm');
await fs.rename(videoPath, finalVideo);

const readme = `# AITrailblazer AI Agent Publishing Three-Minute Demo

Recorded: ${new Date().toISOString()}
URL: ${url}
Viewport: ${width}x${height}
Duration target: ${durationMs} ms
Landing hold: ${url.endsWith('/demo') || url.endsWith('/demo/') ? landingHoldMs : 0} ms

This recording starts on the landing page, then opens the browser demo viewport. The flow demonstrates archive brief generation, TripCode/River resolution, runtime proof packaging, rendered evidence output, second-turn session memory, and usage visibility.

Files:
- 000-start.png
- 010-demo-start.png
- 999-end.png
- AITrailblazer_Publishing_Three_Minute_Demo.webm
- FINAL_TRANSCRIPT.txt
- FINAL_TRANSCRIPT.html
`;
await fs.writeFile(path.join(outDir, 'README.md'), readme);

const transcript = `AITrailblazer AI Agent Publishing - Final Demo Transcript

0:00-0:18 - Landing page and project frame
This project turns a publication archive into agent memory. The opening screen describes the core loop: published research becomes archive context, TripCode identity, River memory, MongoDB-backed state, Gemini-ready synthesis, and a rendered user packet. The checklist shows the submission is a functional running application with Google Cloud runtime surfaces, protected demo access, cost-aware usage, repeatable browser actions, and a MongoDB-oriented memory model.

0:18-0:35 - Start the browser proof
The demo opens the HUT publishing proof from the browser. The viewer sees a controlled workflow rather than an unconstrained chat box. The same deployed service can be tested manually or replayed for the video.

0:35-0:58 - Architecture and workflow overview
The screen explains the operating model. A publication archive is converted into structured agent memory. The Cloud Run agent receives the request, calls the publishing workflow, and returns a packet that can be reused by another agent or reader.

0:58-1:20 - Archive brief
The first runtime step generates an archive brief. Raw JSON appears first as proof of the actual route response. This shows the project is not only a visual mockup; it returns machine-readable fields that can drive follow-up workflows.

1:20-1:42 - TripCode and River resolution
The TripCode HUT-RIVER-001 becomes a stable resolver handle. The response connects article memory, River context, MongoDB fit, and monitor-next actions. This is the key product idea: published content becomes addressable agent memory.

1:42-2:05 - Proof package and rendered packet
The demo shows the proof package as JSON, then renders the same response into a polished HTML report. This proves the system can serve both machines and humans from the same bounded response.

2:05-2:35 - Session memory follow-up
The browser runs a second-turn follow-up. The agent preserves context from the previous River so the user does not have to restate the full archive or TripCode chain.

2:35-3:00 - Usage visibility and close
The usage ledger shows cost-aware operation. The final point is that AITrailblazer AI Agent Publishing is a functional Google Cloud browser and API workflow for turning research archives into reusable agent memory, not a static presentation.
`;
await fs.writeFile(path.join(outDir, 'FINAL_TRANSCRIPT.txt'), transcript);

const transcriptHTML = `<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>AITrailblazer Publishing Demo Transcript</title>
<style>
body{margin:0;background:#05070b;color:#f7fbff;font:16px/1.55 Inter,system-ui,sans-serif;padding:32px}
main{max-width:980px;margin:0 auto;border:1px solid rgba(255,255,255,.14);border-radius:22px;background:linear-gradient(180deg,#121c2b,#070c14);padding:28px;box-shadow:0 30px 90px rgba(0,0,0,.45)}
h1{margin:0 0 8px;font-size:34px;line-height:1}.meta{color:#d8b35f;text-transform:uppercase;letter-spacing:.1em;font-weight:900;font-size:12px}
.row{display:grid;grid-template-columns:120px 1fr;gap:18px;border-top:1px solid rgba(255,255,255,.11);padding:16px 0}.time{color:#70d7ff;font-weight:950}.text{color:#dce8f8;font-size:17px}
</style></head><body><main><div class="meta">Final video narration</div><h1>AITrailblazer AI Agent Publishing</h1>
${transcript.split('\n\n').slice(1).map(block => {
  const [first, ...rest] = block.split('\n');
  const parts = first.split(' - ');
  if (parts.length < 2) return `<p class="text">${first}</p>`;
  return `<section class="row"><div class="time">${parts[0]}</div><div class="text"><strong>${parts.slice(1).join(' - ')}</strong><br>${rest.join(' ')}</div></section>`;
}).join('\n')}
</main></body></html>`;
await fs.writeFile(path.join(outDir, 'FINAL_TRANSCRIPT.html'), transcriptHTML);

console.log(finalVideo);
