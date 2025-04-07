# coding: utf8
import sys
import os

from coleta import coleta_pb2 as Coleta

from headers_keys import (CONTRACHEQUE, CONTRACHEQUE_2024,
                          INDENIZACOES, HEADERS)
import number


def parse_employees(fn, chave_coleta, categoria):
    employees = {}
    counter = 1
    for _, row in fn:
        matricula = str(row["MATRICULA"])
        name = row["NOME"]
        function = row["DESCRICAO"]
        location = row["LOTACAO"]

        membro = Coleta.ContraCheque()
        membro.id_contra_cheque = chave_coleta + "/" + str(counter)
        membro.chave_coleta = chave_coleta
        membro.matricula = matricula
        membro.nome = name
        membro.funcao = function
        membro.local_trabalho = location
        membro.tipo = Coleta.ContraCheque.Tipo.Value("MEMBRO")
        membro.ativo = True
        
        membro.remuneracoes.CopyFrom(
            cria_remuneracao(row, categoria)
        )
        
        employees[matricula] = membro
        counter += 1
            
    return employees


def cria_remuneracao(row, categoria):
    remu_array = Coleta.Remuneracoes()
    items = list(HEADERS[categoria].items())
    for i in range(len(items)):
        key, value = items[i][0], items[i][1]
        remuneracao = Coleta.Remuneracao()
        remuneracao.natureza = Coleta.Remuneracao.Natureza.Value("R")
        remuneracao.categoria = categoria
        remuneracao.item = key
        # Caso o valor seja negativo, ele vai transformar em positivo:
        remuneracao.valor = float(abs(number.format_value(row[value])))

        if categoria == CONTRACHEQUE and value in ["COL08", "COL09", "COL10"]:
            remuneracao.valor = remuneracao.valor * (-1)
            remuneracao.natureza = Coleta.Remuneracao.Natureza.Value("D")
        else: 
            remuneracao.tipo_receita = Coleta.Remuneracao.TipoReceita.Value("O")

        if (
            categoria == CONTRACHEQUE
           ) and value in ["COL01"]:
            remuneracao.tipo_receita = Coleta.Remuneracao.TipoReceita.Value("B")

        remu_array.remuneracao.append(remuneracao)

    return remu_array


def update_employees(fn, employees, categoria):
    for _, row in fn:
        matricula = str(row["MATRICULA"])
        if matricula in employees.keys():
            emp = employees[matricula]
            remu = cria_remuneracao(row, categoria)
            emp.remuneracoes.MergeFrom(remu)
            employees[matricula] = emp
    return employees


def parse(data, chave_coleta, month, year):
    employees = {}
    folha = Coleta.FolhaDePagamento()

    # Puts all parsed employees in the big map
    if int(year) > 2024 or (int(year) == 2024 and int(month) >= 6):
        categoria = CONTRACHEQUE_2024
    else:
        categoria = CONTRACHEQUE

    employees.update(parse_employees(data.contracheque, chave_coleta, categoria))

    update_employees(data.indenizatorias, employees, INDENIZACOES)

    for i in employees.values():
        folha.contra_cheque.append(i)
    return folha
