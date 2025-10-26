### Запуск проекта
1. <pre>cp .env.example .env</pre>
2. Проставить значения в .env
3. <pre>docker compose up --build -d</pre>
4. <pre>make migrate #для миграций</pre>
5. <pre>make dev #(для запуска вне docker но тогда нужно поменять POSTRES_HOST на localhost)</pre>


### Swagger UI
<pre>
/swagger/index.html
</pre>