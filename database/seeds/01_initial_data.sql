	-- ===========================================================
	-- Nexia Health
	-- 01_initial_data.sql
	-- Dados iniciais do sistema
	-- ===========================================================


	-- ###########################################################

	-- ESPECIALIDADES
	-- ###########################################################

	INSERT INTO especialidade (nome, descricao)
	VALUES
	    ('Clínico Geral', 'Atendimento médico geral e avaliação inicial do paciente.'),
	    ('Pediatria', 'Atendimento médico especializado para crianças e adolescentes.'),
	    ('Cardiologia', 'Atendimento especializado em doenças do coração e sistema circulatório.'),
	    ('Ortopedia', 'Atendimento especializado em ossos, músculos e articulações.'),
	    ('Ginecologia', 'Atendimento especializado na saúde da mulher.'),
	    ('Dermatologia', 'Atendimento especializado em doenças da pele, cabelos e unhas.'),
	    ('Oftalmologia', 'Atendimento especializado na saúde dos olhos e da visão.'),
	    ('Neurologia', 'Atendimento especializado em doenças do sistema nervoso.');


	-- ###########################################################
	-- SETORES
	-- ###########################################################


	INSERT INTO setor (nome, descricao)
	VALUES
	    ('Recepção', 'Setor responsável pelo cadastro e atendimento inicial dos pacientes.'),
	    ('Triagem', 'Setor responsável pela avaliação e classificação de risco dos pacientes.'),
	    ('Consultório Médico', 'Setor destinado aos atendimentos médicos.'),
	    ('Enfermagem', 'Setor destinado aos procedimentos e atendimentos de enfermagem.'),
	    ('Laboratório', 'Setor destinado à realização e processamento de exames laboratoriais.'),
	    ('Farmácia', 'Setor destinado à dispensação de medicamentos.'),
	    ('Observação', 'Setor destinado ao acompanhamento de pacientes em observação.');