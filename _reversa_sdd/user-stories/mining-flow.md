# User Story: Solo Mining and View Cycle

> **Nível de Documentação:** COMPLETO  
> **Gerado em:** 2026-05-29

---

## 1. Description
As a Bitcoin enthusiast (Operator),  
I want to run the NerdTUI miner on my local terminal and view real-time statistics  
So that I can monitor hashing performance, shares found, and CPU usage easily without cluttering my system resources.

---

## 2. Key User Scenario: Mining Cycle and View Rotation

### Happy Path (Stratum Mode)
1. **Startup**: The user executes `bin/nerdtui` on a standard UNIX terminal.
2. **Dashboard**: The program launches inside an Alternate Terminal Buffer (`AltScreen`), rendering the **Dashboard Screen** by default.
3. **Loop Hashing**: The background `MinerWorker` registers connection state as `connected` (represented by a green bullet `●` on the StatusBar) and commences double-hashing `SHA256d` loops in 50k batches.
4. **Metrics update**: Every 1s, the UI hashrate sparkline shifts, rendering the current HPS in a scaled formatted layout (e.g. `"12.4 KH/s"`).
5. **View Rotation**: The user presses `tab` key:
   - The view shifts to the **Clock Screen**, showing a large ASCII representation of the current hour.
   - The user presses `tab` key again:
   - The view shifts to the **Global Stats Screen**, showing blocks height and estimated global diff.
6. **CPUTarget Throttling**: The user observes high temperature and presses `-` key:
   - Target is adjusted down by `0.05` ($5\%$).
   - The background worker catches the channel updates and slows down the loop.
   - CPU usage stabilizes and resizes actual metrics.
7. **Exit**: The user presses `q` key. The program gracefully restores original terminal settings and exits.

### Mock Mode Scenario
1. **Startup**: The user executes `bin/nerdtui --mock --cpu 0.3`.
2. **Dashboard**: The program starts up in mock mode without requesting a Bitcoin address, generating mock hashes.
3. **Throttling**: The loop targets exactly 30% CPU usage. Uptime counting behaves exactly as expected.
